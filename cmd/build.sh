#!/bin/bash

# Arc Language Compiler and Test Suite Build Script

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Binaries
ARC_BINARY="arc"
TEST_RUNNER="test_runner"

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}→ $1${NC}"
}

# Build Arc compiler
build_compiler() {
    print_info "Building Arc compiler..."
    GOPROXY=direct go install github.com/arc-language/arc-lang/cmd/arc@latest
    if go build -o "$ARC_BINARY" main.go; then
        print_success "Arc compiler built successfully"
        return 0
    else
        print_error "Failed to build Arc compiler"
        return 1
    fi
}

# Build test runner
build_tests() {
    print_info "Building test runner..."
    
    # Find all test files
    TEST_FILES=$(ls test_*.go 2>/dev/null)
    
    if [ -z "$TEST_FILES" ]; then
        print_error "No test files found (test_*.go)"
        return 1
    fi
    
    if go build -o "$TEST_RUNNER" $TEST_FILES; then
        print_success "Test runner built successfully"
        return 0
    else
        print_error "Failed to build test runner"
        return 1
    fi
}

# Run tests
run_tests() {
    if [ ! -f "./$TEST_RUNNER" ]; then
        print_error "Test runner not found. Building..."
        build_tests || return 1
    fi
    
    if [ ! -f "./$ARC_BINARY" ]; then
        print_error "Arc compiler not found. Building..."
        build_compiler || return 1
    fi
    
    print_info "Running Arc compiler tests..."
    ./"$TEST_RUNNER"
}

# Clean build artifacts
clean() {
    print_info "Cleaning build artifacts..."
    rm -f "$ARC_BINARY" "$TEST_RUNNER"
    rm -f tests.log
    rm -rf /tmp/arc_test_env
    print_success "Clean complete"
}

# Deep clean
clean_all() {
    clean
    print_info "Removing all generated files..."
    rm -f *.o *.out
    print_success "Deep clean complete"
}

# Show help
show_help() {
    echo "Arc Language Compiler - Build Script"
    echo ""
    echo "Usage: ./build.sh [command]"
    echo ""
    echo "Commands:"
    echo "  compiler    Build only the Arc compiler"
    echo "  tests       Build only the test runner"
    echo "  build       Build both compiler and test runner"
    echo "  test        Build and run all tests"
    echo "  run-tests   Run tests without rebuilding"
    echo "  clean       Remove build artifacts"
    echo "  clean-all   Remove all generated files"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./build.sh build"
    echo "  ./build.sh test"
    echo "  ./build.sh clean"
}

# Main script logic
case "${1:-build}" in
    compiler)
        build_compiler
        ;;
    tests)
        build_tests
        ;;
    build)
        build_compiler && build_tests
        ;;
    test)
        build_compiler && build_tests && run_tests
        ;;
    run-tests)
        run_tests
        ;;
    clean)
        clean
        ;;
    clean-all)
        clean_all
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac