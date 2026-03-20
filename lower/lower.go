// lower/lower.go
package lower

import (
	"github.com/arc-language/arc-lang/ast"
)

// Lowerer manages the sequential transformation passes over the AST.
// Passes are ordered carefully: async must run before defer (it generates
// new function bodies that may contain defers), and defer before ARC
// (deferred calls must be visible as ExprStmts before ARC scans them).
type Lowerer struct {
	File *ast.File
}

func NewLowerer(file *ast.File) *Lowerer {
	return &Lowerer{File: file}
}

// Apply runs all lowering passes in the correct order.
func (l *Lowerer) Apply() {
	// Pass 1 — Async: splits async funcs into packet struct + thread entry + wrapper.
	// Must run first because it manufactures new FuncDecls that later passes must see.
	lowerAsync(l.File)

	// Pass 2 — Defer: rewrites defer stmts into explicit calls at every exit point.
	// Must run before ARC so the injected ExprStmts are visible to the lifetime pass.
	lowerDefer(l.File)

	// Pass 3 — ARC: injects decref() calls at every scope exit for 'var' bindings.
	lowerARC(l.File)
}