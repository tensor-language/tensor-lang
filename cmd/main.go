package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/arc-language/arc-lang/backend/backend"
	backendelf "github.com/arc-language/arc-lang/linker/elf"
	"github.com/arc-language/arc-lang/codegen"
	"github.com/arc-language/arc-lang/frontend"
	"github.com/arc-language/arc-lang/lower"
	"github.com/arc-language/arc-lang/syntax"
)

// multiFlag accumulates repeated flags, e.g. -L /usr/lib -L /lib
type multiFlag []string

func (f *multiFlag) String() string { return strings.Join(*f, ", ") }
func (f *multiFlag) Set(v string) error {
	*f = append(*f, v)
	return nil
}

// emitMode controls what the compiler produces
type emitMode int

const (
	emitIR  emitMode = iota // ir  → IR text
	emitObj                 // obj → relocatable .o
	emitExe                 // exe → static ELF executable
	emitBin                 // bin → dynamic ELF executable
)

const usage = `Arc compiler

Usage:
  arc build <file> [flags]   Compile and link a dynamic binary (default)
  arc ir     <file> [flags]  Emit IR text to stdout or -o file
  arc obj    <file> [flags]  Emit relocatable object (.o)
  arc exe    <file> [flags]  Emit static ELF executable

Flags:
  -o <path>        Output file (default derived from source name)
  -L <dir>         Library search path (repeatable)
  -l <name>        Link against shared library, e.g. -l c (repeatable)
  -entry <sym>     Entry point symbol for dynamic binaries (default: _start)
  -debug-ast       Print the lowered AST before codegen
  -print-ir        Print the generated IR to stderr during any build mode
  -dump-ir <path>  Write the generated IR to a file during any build mode
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	// First arg is the subcommand
	subcommand := os.Args[1]

	var mode emitMode
	switch subcommand {
	case "build":
		mode = emitBin
	case "ir":
		mode = emitIR
	case "obj", "object":
		mode = emitObj
	case "exe", "exec", "static":
		mode = emitExe
	case "help", "-h", "--help":
		fmt.Fprint(os.Stderr, usage)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "error: unknown subcommand '%s'\n\n", subcommand)
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	// Second arg is the source file (positional)
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "error: missing source file\n\n")
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	sourceFile := os.Args[2]

	// Parse the remaining flags after the source file
	fs := flag.NewFlagSet("arc", flag.ExitOnError)

	outputFile := fs.String("o", "", "Output file")
	entryPoint := fs.String("entry", "_start", "Entry point symbol (bin mode only)")
	debugAST   := fs.Bool("debug-ast", false, "Print the lowered AST before codegen")
	printIR    := fs.Bool("print-ir", false, "Print generated IR to stderr")
	dumpIR     := fs.String("dump-ir", "", "Write generated IR to this file")

	var libPaths multiFlag
	var linkLibs multiFlag
	fs.Var(&libPaths, "L", "Library search path (repeatable)")
	fs.Var(&linkLibs, "l", "Link against shared library (repeatable)")

	if err := fs.Parse(os.Args[3:]); err != nil {
		fatalf("%v", err)
	}

	// Default output path
	ext        := filepath.Ext(sourceFile)
	moduleName := strings.TrimSuffix(filepath.Base(sourceFile), ext)

	if *outputFile == "" {
		switch mode {
		case emitIR:
			// stdout, leave empty
		case emitObj:
			*outputFile = moduleName + ".o"
		case emitExe, emitBin:
			*outputFile = moduleName
		}
	}

	// Read source
	code, err := os.ReadFile(sourceFile)
	if err != nil {
		fatalf("cannot read '%s': %v", sourceFile, err)
	}

	// Phase 1: Parse
	parseResult := syntax.Parse(string(code))
	if len(parseResult.Errors) > 0 {
		fmt.Fprintln(os.Stderr, "syntax errors:")
		for _, e := range parseResult.Errors {
			fmt.Fprintln(os.Stderr, " ", e)
		}
		os.Exit(1)
	}

	// Phase 2: Semantic analysis
	astFile  := frontend.Translate(parseResult.Root)
	analyzer := frontend.NewAnalyzer()
	if err := analyzer.Analyze(astFile); err != nil {
		fatalf("semantic error: %v", err)
	}

	// Phase 3: Lowering
	lower.NewLowerer(astFile).Apply()

	if *debugAST {
		fmt.Fprintln(os.Stderr, "── lowered AST ──────────────────────────────")
		fmt.Fprintf(os.Stderr, "%+v\n", astFile)
		fmt.Fprintln(os.Stderr, "─────────────────────────────────────────────")
	}

	// Phase 4: IR generation
	gen      := codegen.New(moduleName)
	irModule, err := gen.Generate(astFile)
	if err != nil {
		fatalf("codegen error: %v", err)
	}

	// IR debug output — runs for every mode when requested
	irText := irModule.String()

	if *printIR {
		fmt.Fprintln(os.Stderr, "── generated IR ─────────────────────────────")
		fmt.Fprint(os.Stderr, irText)
		fmt.Fprintln(os.Stderr, "─────────────────────────────────────────────")
	}

	if *dumpIR != "" {
		if err := os.WriteFile(*dumpIR, []byte(irText), 0o644); err != nil {
			fatalf("cannot write IR dump to '%s': %v", *dumpIR, err)
		}
		log.Printf("IR dumped  →  %s", *dumpIR)
	}

	// Phase 5: Backend
	switch mode {

	case emitIR:
		if *outputFile != "" {
			if err := os.WriteFile(*outputFile, []byte(irText), 0o644); err != nil {
				fatalf("cannot write '%s': %v", *outputFile, err)
			}
			log.Printf("compiled '%s'  →  %s  (IR)", sourceFile, *outputFile)
		} else {
			fmt.Print(irText)
		}

	case emitObj:
		if err := backend.Generate(irModule); err != nil {
			fatalf("backend error: %v", err)
		}
		objBytes, err := backend.GenerateObject(irModule)
		if err != nil {
			fatalf("object generation failed: %v", err)
		}
		if err := writeFile(*outputFile, objBytes, 0o644); err != nil {
			fatalf("cannot write '%s': %v", *outputFile, err)
		}
		log.Printf("compiled '%s'  →  %s  (%d bytes)", sourceFile, *outputFile, len(objBytes))

	case emitExe:
		if err := backend.Generate(irModule); err != nil {
			fatalf("backend error: %v", err)
		}
		exeBytes, err := backend.GenerateExecutable(irModule)
		if err != nil {
			fatalf("executable generation failed: %v", err)
		}
		if err := writeFile(*outputFile, exeBytes, 0o755); err != nil {
			fatalf("cannot write '%s': %v", *outputFile, err)
		}
		log.Printf("compiled '%s'  →  %s  (%d bytes, static)", sourceFile, *outputFile, len(exeBytes))

	case emitBin:
		if err := backend.Generate(irModule); err != nil {
			fatalf("backend error: %v", err)
		}
		objBytes, err := backend.GenerateObject(irModule)
		if err != nil {
			fatalf("object generation failed: %v", err)
		}

		linker := backendelf.NewLinker(backendelf.Config{
			Entry:    *entryPoint,
			BaseAddr: 0x400000,
		})
		if err := linker.AddObject(moduleName+".o", objBytes); err != nil {
			fatalf("linker: failed to add object: %v", err)
		}
		if err := resolveSharedLibs(linker, linkLibs, libPaths); err != nil {
			fatalf("linker: %v", err)
		}
		if err := linker.Link(*outputFile); err != nil {
			fatalf("link failed: %v", err)
		}
		log.Printf("compiled '%s'  →  %s  (dynamic)", sourceFile, *outputFile)
	}
}

// resolveSharedLibs finds and loads each -l library into the linker.
func resolveSharedLibs(linker *backendelf.Linker, libs, paths multiFlag) error {
	searchDirs := append([]string(paths), systemLibDirs()...)

	for _, lib := range libs {
		candidates := libCandidates(lib)
		found := false

		for _, dir := range searchDirs {
			for _, cand := range candidates {
				full := filepath.Join(dir, cand)
				data, err := os.ReadFile(full)
				if err != nil {
					continue
				}
				// Skip linker scripts and other non-ELF files
				if !isELF(data) {
					continue
				}
				if err := linker.AddSharedLib(full, data); err != nil {
					log.Printf("warning: skipping '%s': %v", full, err)
					continue
				}
				log.Printf("linked shared library  %s", full)
				found = true
				break
			}
			if found {
				break
			}
		}

		if !found {
			return fmt.Errorf("shared library 'lib%s' not found in search paths", lib)
		}
	}
	return nil
}

// isELF returns true if data begins with the ELF magic number.
func isELF(data []byte) bool {
	return len(data) >= 4 &&
		data[0] == 0x7f &&
		data[1] == 'E' &&
		data[2] == 'L' &&
		data[3] == 'F'
}

// libCandidates returns filenames to try for a -l<name> flag.
// Versioned names are tried before the bare .so to avoid linker scripts.
func libCandidates(name string) []string {
	if strings.HasPrefix(name, "lib") &&
		(strings.HasSuffix(name, ".so") || strings.Contains(name, ".so.")) {
		return []string{name}
	}
	return []string{
		"lib" + name + ".so.6",
		"lib" + name + ".so.2",
		"lib" + name + ".so.5",
		"lib" + name + ".so.1",
		"lib" + name + ".so",
	}
}

// systemLibDirs returns the standard shared library search paths on Linux.
func systemLibDirs() []string {
	return []string{
		"/lib/x86_64-linux-gnu",
		"/usr/lib/x86_64-linux-gnu",
		"/lib64",
		"/usr/lib64",
		"/lib",
		"/usr/lib",
	}
}

// writeFile writes data to path and sets the given permissions.
func writeFile(path string, data []byte, perm os.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return err
	}
	return os.Chmod(path, perm)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}