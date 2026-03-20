// lower/walk.go
package lower

import "github.com/arc-language/arc-lang/ast"

// walkDecls iterates the top-level declarations in a file and calls fn for
// each FuncDecl body. All three lowering passes share this entry point.
func walkDecls(file *ast.File, fn func(*ast.FuncDecl)) {
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			fn(f)
		}
	}
}

// walkStmts applies a visitor to every statement reachable from b, providing
// a chance to rewrite the statement slice in place. The visitor returns the
// replacement slice for each block it processes; the walk recurses into nested
// blocks after the visitor returns so that new statements are also visited.
//
// visitor receives the current block and its parent scope depth, and must
// return the replacement statement list (it may be the same slice).
func walkBlock(b *ast.BlockStmt, visitor func(stmts []ast.Stmt) []ast.Stmt) {
	b.List = visitor(b.List)
	for _, stmt := range b.List {
		walkStmtChildren(stmt, visitor)
	}
}

// walkStmtChildren recurses into all compound statements that contain blocks
// or statement lists, applying visitor to each inner block it finds.
func walkStmtChildren(stmt ast.Stmt, visitor func(stmts []ast.Stmt) []ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		walkBlock(s, visitor)

	case *ast.IfStmt:
		walkBlock(s.Body, visitor)
		if s.Else != nil {
			walkStmtChildren(s.Else, visitor)
		}

	case *ast.ForStmt:
		if s.Init != nil {
			walkStmtChildren(s.Init, visitor)
		}
		walkBlock(s.Body, visitor)

	case *ast.ForInStmt:
		walkBlock(s.Body, visitor)

	case *ast.SwitchStmt:
		for _, c := range s.Cases {
			c.Body = visitor(c.Body)
			for _, st := range c.Body {
				walkStmtChildren(st, visitor)
			}
		}
		if s.Default != nil {
			s.Default = visitor(s.Default)
			for _, st := range s.Default {
				walkStmtChildren(st, visitor)
			}
		}
	}
}

// walkExprs visits every expression reachable from root, calling fn on each
// one. fn receives the expression and returns its replacement (may be the same
// pointer). This is used by the await-rewriting step.
func walkExprs(root ast.Node, fn func(ast.Expr) ast.Expr) {
	switch n := root.(type) {
	case *ast.File:
		for _, d := range n.Decls {
			walkExprs(d, fn)
		}
	case *ast.FuncDecl:
		if n.Body != nil {
			walkExprsBlock(n.Body, fn)
		}
	case *ast.DeinitDecl:
		if n.Body != nil {
			walkExprsBlock(n.Body, fn)
		}
	}
}

func walkExprsBlock(b *ast.BlockStmt, fn func(ast.Expr) ast.Expr) {
	for i, stmt := range b.List {
		b.List[i] = walkExprsStmt(stmt, fn)
	}
}

func walkExprsStmt(stmt ast.Stmt, fn func(ast.Expr) ast.Expr) ast.Stmt {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		s.X = walkExprsExpr(s.X, fn)
	case *ast.DeclStmt:
		if v, ok := s.Decl.(*ast.VarDecl); ok && v.Value != nil {
			v.Value = walkExprsExpr(v.Value, fn)
		}
	case *ast.AssignStmt:
		s.Target = walkExprsExpr(s.Target, fn)
		if s.Value != nil {
			s.Value = walkExprsExpr(s.Value, fn)
		}
	case *ast.ReturnStmt:
		for i, r := range s.Results {
			s.Results[i] = walkExprsExpr(r, fn)
		}
	case *ast.DeferStmt:
		s.Call = walkExprsExpr(s.Call, fn)
	case *ast.IfStmt:
		s.Cond = walkExprsExpr(s.Cond, fn)
		walkExprsBlock(s.Body, fn)
		if s.Else != nil {
			s.Else = walkExprsStmt(s.Else, fn)
		}
	case *ast.ForStmt:
		if s.Init != nil {
			s.Init = walkExprsStmt(s.Init, fn)
		}
		if s.Cond != nil {
			s.Cond = walkExprsExpr(s.Cond, fn)
		}
		if s.Post != nil {
			s.Post = walkExprsStmt(s.Post, fn)
		}
		walkExprsBlock(s.Body, fn)
	case *ast.ForInStmt:
		s.Iter = walkExprsExpr(s.Iter, fn)
		walkExprsBlock(s.Body, fn)
	case *ast.SwitchStmt:
		s.Tag = walkExprsExpr(s.Tag, fn)
		for _, c := range s.Cases {
			for i, v := range c.Values {
				c.Values[i] = walkExprsExpr(v, fn)
			}
			for i, st := range c.Body {
				c.Body[i] = walkExprsStmt(st, fn)
			}
		}
		for i, st := range s.Default {
			s.Default[i] = walkExprsStmt(st, fn)
		}
	case *ast.BlockStmt:
		walkExprsBlock(s, fn)
	}
	return stmt
}

func walkExprsExpr(expr ast.Expr, fn func(ast.Expr) ast.Expr) ast.Expr {
	if expr == nil {
		return nil
	}
	// First recurse into children so fn sees the innermost nodes first
	// (post-order), which matches what both await-rewriting and ARC need.
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		e.Left = walkExprsExpr(e.Left, fn)
		e.Right = walkExprsExpr(e.Right, fn)
	case *ast.UnaryExpr:
		e.X = walkExprsExpr(e.X, fn)
	case *ast.CallExpr:
		e.Fun = walkExprsExpr(e.Fun, fn)
		for i, arg := range e.Args {
			e.Args[i] = walkExprsExpr(arg, fn)
		}
	case *ast.SelectorExpr:
		e.X = walkExprsExpr(e.X, fn)
	case *ast.IndexExpr:
		e.X = walkExprsExpr(e.X, fn)
		e.Index = walkExprsExpr(e.Index, fn)
	case *ast.SliceExpr:
		e.X = walkExprsExpr(e.X, fn)
		e.Low = walkExprsExpr(e.Low, fn)
		e.High = walkExprsExpr(e.High, fn)
	case *ast.RangeExpr:
		e.Low = walkExprsExpr(e.Low, fn)
		e.High = walkExprsExpr(e.High, fn)
	case *ast.AwaitExpr:
		e.X = walkExprsExpr(e.X, fn)
	case *ast.CompositeLit:
		for i, f := range e.Fields {
			e.Fields[i] = walkExprsExpr(f, fn)
		}
	case *ast.KeyValueExpr:
		e.Value = walkExprsExpr(e.Value, fn)
	case *ast.TupleLit:
		for i, el := range e.Elems {
			e.Elems[i] = walkExprsExpr(el, fn)
		}
	case *ast.LambdaExpr:
		walkExprsBlock(e.Body, fn)
	case *ast.ProcessExpr:
		for i, arg := range e.Args {
			e.Args[i] = walkExprsExpr(arg, fn)
		}
		walkExprsBlock(e.Body, fn)
	case *ast.NewExpr:
		if e.Init != nil {
			e.Init = walkExprsExpr(e.Init, fn).(*ast.CompositeLit)
		}
	case *ast.NewArrayExpr:
		e.Len = walkExprsExpr(e.Len, fn)
	case *ast.DeleteExpr:
		e.X = walkExprsExpr(e.X, fn)
	case *ast.CastExpr:
		e.X = walkExprsExpr(e.X, fn)
	}
	return fn(expr)
}