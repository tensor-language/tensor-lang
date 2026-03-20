// lower/arc.go
package lower

import "github.com/arc-language/arc-lang/ast"

// lowerARC injects decref() calls at every scope exit for reference-counted
// ('var') bindings. The strategy is conservative:
//
//   - 'let' bindings are not reference-counted (they are immutable value types
//     or borrows in Arc's ownership model).
//   - 'var' bindings are heap-allocated and ref-counted; we decref at every
//     point where they leave scope (return, end-of-block).
//   - Each nested block (if-body, loop-body, etc.) maintains its own var list
//     and drains it when the block exits, restoring the parent list afterward.
//
// decref is treated as a compiler intrinsic. In later codegen it resolves to
// the runtime symbol that decrements the reference count and frees if zero.
func lowerARC(file *ast.File) {
	walkDecls(file, func(fn *ast.FuncDecl) {
		if fn.Body != nil {
			ai := &arcInjector{}
			ai.processBlock(fn.Body)
		}
	})
}

type arcInjector struct {
	// vars is the list of 'var' declarations active in the current block.
	vars []*ast.VarDecl
}

// processBlock rewrites b in place, injecting decref calls at exit points.
func (ai *arcInjector) processBlock(b *ast.BlockStmt) {
	// Save and reset per-block var list.
	saved := ai.vars
	ai.vars = nil

	var out []ast.Stmt

	for _, stmt := range b.List {
		switch s := stmt.(type) {

		case *ast.DeclStmt:
			out = append(out, s)
			// Track 'var' (IsRef) bindings only; 'let' bindings are not RC.
			if v, ok := s.Decl.(*ast.VarDecl); ok && v.IsRef {
				ai.vars = append(ai.vars, v)
			}

		case *ast.ReturnStmt:
			// Decref all active vars before returning.
			ai.emitDecrefs(&out)
			out = append(out, s)

		case *ast.IfStmt:
			ai.processIf(s)
			out = append(out, s)

		case *ast.ForStmt:
			if s.Body != nil {
				ai.processBlock(s.Body)
			}
			out = append(out, s)

		case *ast.ForInStmt:
			if s.Body != nil {
				ai.processBlock(s.Body)
			}
			out = append(out, s)

		case *ast.SwitchStmt:
			ai.processSwitch(s)
			out = append(out, s)

		case *ast.BlockStmt:
			ai.processBlock(s)
			out = append(out, s)

		default:
			out = append(out, s)
		}
	}

	// Drain at end-of-block (covers fall-through paths and void functions).
	ai.emitDecrefs(&out)

	b.List = out
	ai.vars = saved
}

// processIf recurses into both branches. Each branch inherits the parent var
// list (as a snapshot) so vars declared before the if are decref'd on any
// early return inside the branch.
func (ai *arcInjector) processIf(s *ast.IfStmt) {
	ai.processBlock(s.Body)
	if s.Else != nil {
		switch e := s.Else.(type) {
		case *ast.BlockStmt:
			ai.processBlock(e)
		case *ast.IfStmt:
			ai.processIf(e)
		}
	}
}

// processSwitch handles each case body independently.
func (ai *arcInjector) processSwitch(s *ast.SwitchStmt) {
	for _, c := range s.Cases {
		saved := ai.vars
		ai.vars = nil

		var out []ast.Stmt
		for _, stmt := range c.Body {
			if ret, ok := stmt.(*ast.ReturnStmt); ok {
				ai.emitDecrefs(&out)
				out = append(out, ret)
			} else {
				if d, ok := stmt.(*ast.DeclStmt); ok {
					if v, ok2 := d.Decl.(*ast.VarDecl); ok2 && v.IsRef {
						ai.vars = append(ai.vars, v)
					}
				}
				out = append(out, stmt)
			}
		}
		ai.emitDecrefs(&out)
		c.Body = out
		ai.vars = saved
	}

	if s.Default != nil {
		saved := ai.vars
		ai.vars = nil

		var out []ast.Stmt
		for _, stmt := range s.Default {
			if ret, ok := stmt.(*ast.ReturnStmt); ok {
				ai.emitDecrefs(&out)
				out = append(out, ret)
			} else {
				if d, ok := stmt.(*ast.DeclStmt); ok {
					if v, ok2 := d.Decl.(*ast.VarDecl); ok2 && v.IsRef {
						ai.vars = append(ai.vars, v)
					}
				}
				out = append(out, stmt)
			}
		}
		ai.emitDecrefs(&out)
		s.Default = out
		ai.vars = saved
	}
}

// emitDecrefs appends a decref(v) ExprStmt for every tracked var in LIFO order.
func (ai *arcInjector) emitDecrefs(out *[]ast.Stmt) {
	for i := len(ai.vars) - 1; i >= 0; i-- {
		v := ai.vars[i]
		*out = append(*out, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "decref"},
				Args: []ast.Expr{&ast.Ident{Name: v.Name}},
			},
		})
	}
}