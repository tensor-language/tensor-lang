package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const TEST_TIMEOUT = 10 * time.Second
const LIB_PATH = "/usr/lib/x86_64-linux-gnu"
const TEST_DIR = "../tests"

type TestCase struct {
	Name     string
	FilePath string
	Expected string
}

func main() {
	logFile, err := os.Create("/tmp/tests.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	arcBinary := "./arc"
	if _, err := os.Stat(arcBinary); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Error: './arc' binary not found. Please build the compiler first.")
		os.Exit(1)
	}

	tempDir := filepath.Join(os.TempDir(), "arc_test_env")
	os.RemoveAll(tempDir)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	// Discover all .ax files recursively under TEST_DIR
	tests, err := discoverTests(TEST_DIR)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering tests: %v\n", err)
		os.Exit(1)
	}

	if len(tests) == 0 {
		fmt.Println("No test files found.")
		os.Exit(0)
	}

	// Sort alphabetically by file path
	sort.Slice(tests, func(i, j int) bool {
		return tests[i].FilePath < tests[j].FilePath
	})

	fmt.Printf("Arc Compiler Test Suite - Running %d tests...\n", len(tests))
	fmt.Println(strings.Repeat("=", 60))

	passed := 0
	failed := 0
	startTotal := time.Now()

	for i, tc := range tests {
		prefix := fmt.Sprintf("[%d/%d] %-40s ", i+1, len(tests), tc.Name)
		fmt.Print(prefix)
		logFile.WriteString(fmt.Sprintf("\n--- Test: %s (%s) ---\n", tc.Name, tc.FilePath))

		start := time.Now()
		output, err := runTest(tc, arcBinary, tempDir, logFile)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("❌ FAIL (%.3fs)\n", duration.Seconds())
			fmt.Printf("      Error: %v\n", err)
			logFile.WriteString(fmt.Sprintf("Result: FAIL\nError: %v\n", err))
			failed++
		} else {
			fmt.Printf("✅ PASS (%.3fs)\n", duration.Seconds())
			logFile.WriteString("Result: PASS\n")
			passed++
		}

		if output != "" {
			fmt.Printf("      Output: %s\n", strings.TrimSpace(output))
		}
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Passed: %d | Failed: %d | Total Time: %.2fs\n", passed, failed, time.Since(startTotal).Seconds())

	if failed > 0 {
		os.Exit(1)
	}
}

// discoverTests walks TEST_DIR recursively, finds all .ax files,
// and parses the //@test(expect = "...") header from each.
func discoverTests(root string) ([]TestCase, error) {
	var tests []TestCase

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".ax" {
			return nil
		}

		expected, err := parseExpectHeader(path)
		if err != nil {
			// Skip files that don't have a valid header rather than failing
			fmt.Fprintf(os.Stderr, "  Warning: skipping %s: %v\n", path, err)
			return nil
		}

		tests = append(tests, TestCase{
			Name:     strings.TrimSuffix(filepath.Base(path), ".ax"),
			FilePath: path,
			Expected: expected,
		})
		return nil
	})

	return tests, err
}

// unescapeString converts escape sequences in the expected string
func unescapeString(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\t", "\t")
	s = strings.ReplaceAll(s, "\\r", "\r")
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "\\\\", "\\")
	return s
}

// parseExpectHeader reads the first line of the file looking for:
//
//	//@test(expect = "...")
func parseExpectHeader(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return "", fmt.Errorf("empty file")
	}

	line := strings.TrimSpace(scanner.Text())

	// Must start with //@test(expect =
	if !strings.HasPrefix(line, "//@test(") {
		return "", fmt.Errorf("no //@test header found")
	}

	// Extract the value between the quotes after expect =
	eqIdx := strings.Index(line, "expect")
	if eqIdx == -1 {
		return "", fmt.Errorf("no expect field in @test header")
	}

	rest := line[eqIdx:]
	openQuote := strings.Index(rest, "\"")
	if openQuote == -1 {
		return "", fmt.Errorf("missing opening quote in expect")
	}
	closeQuote := strings.Index(rest[openQuote+1:], "\"")
	if closeQuote == -1 {
		return "", fmt.Errorf("missing closing quote in expect")
	}

	expected := rest[openQuote+1 : openQuote+1+closeQuote]
	return unescapeString(expected), nil
}

func runTest(tc TestCase, arcBinary, tempDir string, log *os.File) (string, error) {
	exePath := filepath.Join(tempDir, tc.Name)

	// Compile
	ctx, cancel := context.WithTimeout(context.Background(), TEST_TIMEOUT)
	defer cancel()

	cmdBuild := exec.CommandContext(ctx, arcBinary,
		"build", tc.FilePath,
		"-o", exePath,
		"-L", LIB_PATH,
		"-l", "c",
	)

	outBuild, err := cmdBuild.CombinedOutput()
	log.WriteString("Build Output:\n" + string(outBuild) + "\n")

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("build timeout (exceeded %v)", TEST_TIMEOUT)
	}
	if err != nil {
		return "", fmt.Errorf("build failed: %s", strings.TrimSpace(string(outBuild)))
	}

	// Run
	ctx, cancel = context.WithTimeout(context.Background(), TEST_TIMEOUT)
	defer cancel()

	cmdRun := exec.CommandContext(ctx, exePath)
	outRun, err := cmdRun.CombinedOutput()
	output := string(outRun)
	log.WriteString("Run Output:\n" + output + "\n")

	if ctx.Err() == context.DeadlineExceeded {
		return output, fmt.Errorf("execution timeout (exceeded %v) - possible infinite loop", TEST_TIMEOUT)
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return output, fmt.Errorf("runtime error (exit code %d)\nOutput: %s", exitErr.ExitCode(), output)
		}
		return output, fmt.Errorf("runtime error: %v\nOutput: %s", err, output)
	}

	// Assert
	if !strings.Contains(output, tc.Expected) {
		return output, fmt.Errorf("assertion failed.\n      Expected (quoted): %q\n      Expected (raw): %s\n      Got (quoted): %q\n      Got (raw): %s", 
			tc.Expected, tc.Expected, 
			strings.TrimSpace(output), strings.TrimSpace(output))
	}

	return output, nil
}