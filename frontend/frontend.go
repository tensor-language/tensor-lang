// frontend/frontend.go
package frontend

import (
	"github.com/arc-language/arc-lang/ast"
	"github.com/arc-language/arc-lang/parser"
)

// Result bundles the outputs of the frontend pipeline.
type Result struct {
	File     *ast.File
	Analyzer *Analyzer
}

// Run translates the CST and runs both semantic passes.
// Syntax errors are expected to have been caught by the syntax package already.
func Run(root parser.ICompilationUnitContext) (*Result, error) {
	file := Translate(root)

	analyzer := NewAnalyzer()
	if err := analyzer.Analyze(file); err != nil {
		return &Result{File: file, Analyzer: analyzer}, err
	}
	return &Result{File: file, Analyzer: analyzer}, nil
}