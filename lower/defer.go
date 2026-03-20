// lower/defer.go
package lower

import "github.com/arc-language/arc-lang/ast"

// lowerDefer rewrites defer statements throughout the file.
//
// For every block, we collect deferred expressions as we encounter them, then
// inject them (in LIFO order) before every ReturnStmt and at the end of the
// block. If a defer appears inside a loop body it is scoped to that iteration
// because walkBlock is called recursively, giving each nested block its own
// defer stack.
//
// Limitation: break/continue out of a loop do not currently drain defers. A
// production implementation would use a finally-block or setjmp approach.
func lowerDefer(file *ast.File) {
	walkDecls(file, func(fn *ast.FuncDecl) {
		if fn.Body != nil {
			dl := &deferLowerer{}
			dl.processBlock(fn.Body)
		}
	})
}

type deferLowerer struct {
	// defers is the LIFO stack of expressions for the current block level.
	// Each call to processBlock saves/restores this field so that nesting
	// works correctly.
	defers []ast.Expr
}

// processBlock rewrites a single block in place.
func (dl *deferLowerer) processBlock(b *ast.BlockStmt) {
	// Save the caller's defer stack and start a fresh one for this block.
	saved := dl.defers
	dl.defers = nil

	var out []ast.Stmt

	for _, stmt := range b.List {
		switch s := stmt.(type) {

		case *ast.DeferStmt:
			// Push the deferred expression; do NOT emit the DeferStmt itself.
			dl.defers = append(dl.defers, s.Call)

		case *ast.ReturnStmt:
			// Drain defers in LIFO order before the return.
			dl.emitDefers(&out)
			out = append(out, s)

		case *ast.IfStmt:
			dl.processIf(s)
			out = append(out, s)

		case *ast.ForStmt:
			// Each iteration's defers are scoped to the loop body.
			if s.Body != nil {
				dl.processBlock(s.Body)
			}
			out = append(out, s)

		case *ast.ForInStmt:
			if s.Body != nil {
				dl.processBlock(s.Body)
			}
			out = append(out, s)

		case *ast.SwitchStmt:
			dl.processSwitch(s)
			out = append(out, s)

		case *ast.BlockStmt:
			// Nested bare block (e.g. scoping block) gets its own defer scope.
			dl.processBlock(s)
			out = append(out, s)

		default:
			out = append(out, s)
		}
	}

	// End-of-block drain: covers functions that fall off the end without an
	// explicit return, and any defers not yet consumed by a return above.
	dl.emitDefers(&out)

	b.List = out
	dl.defers = saved
}

// processIf recurses into both branches without merging their defer stacks,
// because each branch has its own exit paths.
func (dl *deferLowerer) processIf(s *ast.IfStmt) {
	dl.processBlock(s.Body)
	if s.Else != nil {
		switch e := s.Else.(type) {
		case *ast.BlockStmt:
			dl.processBlock(e)
		case *ast.IfStmt:
			dl.processIf(e)
		}
	}
}

// processSwitch recurses into every case body independently.
func (dl *deferLowerer) processSwitch(s *ast.SwitchStmt) {
	for _, c := range s.Cases {
		saved := dl.defers
		dl.defers = nil

		var out []ast.Stmt
		for _, stmt := range c.Body {
			if ret, ok := stmt.(*ast.ReturnStmt); ok {
				dl.emitDefers(&out)
				out = append(out, ret)
			} else {
				out = append(out, stmt)
			}
		}
		dl.emitDefers(&out)
		c.Body = out

		dl.defers = saved
	}
	if s.Default != nil {
		saved := dl.defers
		dl.defers = nil

		var out []ast.Stmt
		for _, stmt := range s.Default {
			if ret, ok := stmt.(*ast.ReturnStmt); ok {
				dl.emitDefers(&out)
				out = append(out, ret)
			} else {
				out = append(out, stmt)
			}
		}
		dl.emitDefers(&out)
		s.Default = out

		dl.defers = saved
	}
}

// emitDefers appends ExprStmts for each deferred expression in LIFO order.
// The defers slice is intentionally NOT cleared here: the same stack may be
// drained more than once (once per ReturnStmt, and once at end-of-block),
// which is correct because each return is an independent exit path.
//
// Note: draining does not double-execute — only one exit path is taken at
// runtime. The duplicate injections are dead code for all but one path.
func (dl *deferLowerer) emitDefers(out *[]ast.Stmt) {
	for i := len(dl.defers) - 1; i >= 0; i-- {
		*out = append(*out, &ast.ExprStmt{X: dl.defers[i]})
	}
}