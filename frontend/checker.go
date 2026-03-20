// frontend/checker.go
package frontend

import (
	"fmt"

	"github.com/arc-language/arc-lang/ast"
)

// TypeError carries a positioned semantic error.
type TypeError struct {
	Pos ast.Position
	Msg string
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Pos.Line, e.Pos.Column, e.Msg)
}

// Analyzer runs the two-pass semantic analysis on a translated AST.
type Analyzer struct {
	GlobalScope *Scope
	Errors      []*TypeError
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{GlobalScope: NewScope(nil, ScopeGlobal)}
}

// Analyze runs Pass 1 (symbol discovery) then Pass 2 (type checking).
// Returns the first error encountered, or nil. All errors are also in Errors.
func (a *Analyzer) Analyze(file *ast.File) error {
	a.discoverSymbols(file)
	a.checkFile(file)
	if len(a.Errors) > 0 {
		return a.Errors[0]
	}
	return nil
}

func (a *Analyzer) errorf(pos ast.Position, format string, args ...any) {
	a.Errors = append(a.Errors, &TypeError{Pos: pos, Msg: fmt.Sprintf(format, args...)})
}

// ─── Pass 1: Symbol Discovery ─────────────────────────────────────────────────

func (a *Analyzer) discoverSymbols(file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if prev := a.GlobalScope.Insert(d.Name, &Symbol{Name: d.Name, Kind: "func", Decl: d}); prev != nil {
				a.errorf(d.Start, "func %q redeclared in this block", d.Name)
			}
		case *ast.InterfaceDecl:
			if prev := a.GlobalScope.Insert(d.Name, &Symbol{Name: d.Name, Kind: "type", Decl: d}); prev != nil {
				a.errorf(d.Start, "type %q redeclared in this block", d.Name)
			}
		case *ast.EnumDecl:
			if prev := a.GlobalScope.Insert(d.Name, &Symbol{Name: d.Name, Kind: "enum", Decl: d}); prev != nil {
				a.errorf(d.Start, "enum %q redeclared in this block", d.Name)
			}
			for _, m := range d.Members {
				qual := d.Name + "." + m.Name
				a.GlobalScope.Insert(qual, &Symbol{Name: qual, Kind: "enumMember", Decl: m})
			}
		case *ast.ConstDecl:
			for _, spec := range d.Specs {
				if prev := a.GlobalScope.Insert(spec.Name, &Symbol{Name: spec.Name, Kind: "const", Decl: spec}); prev != nil {
					a.errorf(spec.Start, "const %q redeclared in this block", spec.Name)
				}
			}
		case *ast.VarDecl:
			kind := "let"
			if d.IsRef {
				kind = "var"
			}
			if prev := a.GlobalScope.Insert(d.Name, &Symbol{Name: d.Name, Kind: kind, Decl: d}); prev != nil {
				a.errorf(d.Start, "%s %q redeclared in this block", kind, d.Name)
			}
		case *ast.TypeAliasDecl:
			if prev := a.GlobalScope.Insert(d.Name, &Symbol{Name: d.Name, Kind: "type", Decl: d}); prev != nil {
				a.errorf(d.Start, "type %q redeclared in this block", d.Name)
			}
		case *ast.ExternDecl:
			a.discoverExternSymbols(d, a.GlobalScope)
		}
	}
}

func (a *Analyzer) discoverExternSymbols(d *ast.ExternDecl, scope *Scope) {
	for _, m := range d.Members {
		switch em := m.(type) {
		case *ast.ExternFunc:
			scope.Insert(em.Name, &Symbol{Name: em.Name, Kind: "func", Decl: nil})
		case *ast.ExternNamespace:
			a.discoverExternNamespace(em, em.Name, scope)
		case *ast.ExternClass:
			scope.Insert(em.Name, &Symbol{Name: em.Name, Kind: "type", Decl: nil})
		}
	}
}

func (a *Analyzer) discoverExternNamespace(ns *ast.ExternNamespace, prefix string, scope *Scope) {
	for _, m := range ns.Members {
		switch em := m.(type) {
		case *ast.ExternFunc:
			qual := prefix + "." + em.Name
			scope.Insert(qual, &Symbol{Name: qual, Kind: "func", Decl: nil})
		case *ast.ExternNamespace:
			a.discoverExternNamespace(em, prefix+"."+em.Name, scope)
		case *ast.ExternClass:
			qual := prefix + "." + em.Name
			scope.Insert(qual, &Symbol{Name: qual, Kind: "type", Decl: nil})
		}
	}
}

// ─── Pass 2: Type Checking ─────────────────────────────────────────────────────

func (a *Analyzer) checkFile(file *ast.File) {
	for _, decl := range file.Decls {
		a.checkDecl(decl, a.GlobalScope)
	}
}

func (a *Analyzer) checkDecl(decl ast.Decl, scope *Scope) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		a.checkFuncDecl(d, scope)
	case *ast.DeinitDecl:
		a.checkBlock(d.Body, NewScope(scope, ScopeFunc))
	case *ast.ConstDecl:
		for _, spec := range d.Specs {
			if spec.Value != nil {
				a.checkExpr(spec.Value, scope)
			}
		}
	case *ast.VarDecl:
		if d.Value != nil {
			a.checkExpr(d.Value, scope)
		}
	case *ast.EnumDecl:
		for _, m := range d.Members {
			if m.Value != nil {
				a.checkExpr(m.Value, scope)
			}
		}
	}
}

func (a *Analyzer) checkFuncDecl(d *ast.FuncDecl, parent *Scope) {
	fnScope := NewScope(parent, ScopeFunc)
	if d.Self != nil {
		fnScope.Insert(d.Self.Name, &Symbol{Name: d.Self.Name, Kind: "param", Decl: d, Type: d.Self.Type})
	}
	for _, p := range d.Params {
		if prev := fnScope.Insert(p.Name, &Symbol{Name: p.Name, Kind: "param", Decl: p, Type: p.Type}); prev != nil {
			a.errorf(p.Start, "parameter %q redeclared", p.Name)
		}
	}
	if d.IsAsync && d.IsGpu {
		a.errorf(d.Start, "func %q cannot be both async and gpu", d.Name)
	}
	if d.Body != nil {
		a.checkBlock(d.Body, fnScope)
	}
}

func (a *Analyzer) checkBlock(b *ast.BlockStmt, parent *Scope) {
	scope := NewScope(parent, ScopeBlock)
	for _, stmt := range b.List {
		a.checkStmt(stmt, scope)
	}
}

func (a *Analyzer) checkStmt(stmt ast.Stmt, scope *Scope) {
	switch s := stmt.(type) {
	case *ast.DeclStmt:
		a.checkDeclStmt(s, scope)
	case *ast.AssignStmt:
		a.checkAssignStmt(s, scope)
	case *ast.ExprStmt:
		a.checkExpr(s.X, scope)
	case *ast.ReturnStmt:
		for _, r := range s.Results {
			a.checkExpr(r, scope)
		}
	case *ast.DeferStmt:
		a.checkExpr(s.Call, scope)
	case *ast.IfStmt:
		a.checkExpr(s.Cond, scope)
		a.checkBlock(s.Body, scope)
		if s.Else != nil {
			a.checkStmt(s.Else, scope)
		}
	case *ast.BlockStmt:
		a.checkBlock(s, scope)
	case *ast.ForStmt:
		a.checkForStmt(s, scope)
	case *ast.ForInStmt:
		a.checkExpr(s.Iter, scope)
		loopScope := NewScope(scope, ScopeBlock)
		loopScope.Insert(s.Key, &Symbol{Name: s.Key, Kind: "let"})
		if s.Value != "" {
			loopScope.Insert(s.Value, &Symbol{Name: s.Value, Kind: "let"})
		}
		a.checkBlock(s.Body, loopScope)
	case *ast.SwitchStmt:
		a.checkExpr(s.Tag, scope)
		for _, c := range s.Cases {
			for _, v := range c.Values {
				a.checkExpr(v, scope)
			}
			caseScope := NewScope(scope, ScopeBlock)
			for _, st := range c.Body {
				a.checkStmt(st, caseScope)
			}
		}
		if s.Default != nil {
			defScope := NewScope(scope, ScopeBlock)
			for _, st := range s.Default {
				a.checkStmt(st, defScope)
			}
		}
	}
}

func (a *Analyzer) checkDeclStmt(s *ast.DeclStmt, scope *Scope) {
	switch d := s.Decl.(type) {
	case *ast.VarDecl:
		if d.Type != nil && d.Value != nil {
			a.inferCompositeLitType(d.Value, d.Type)
		}
		if d.Value != nil {
			a.checkExpr(d.Value, scope)
		}
		kind := "let"
		if d.IsRef {
			kind = "var"
		}
		if prev := scope.Insert(d.Name, &Symbol{Name: d.Name, Kind: kind, Decl: d, Type: d.Type}); prev != nil {
			a.errorf(d.Start, "%s %q redeclared in this block", kind, d.Name)
		}
	case *ast.ConstDecl:
		for _, spec := range d.Specs {
			if spec.Type != nil && spec.Value != nil {
				a.inferCompositeLitType(spec.Value, spec.Type)
			}
			if spec.Value != nil {
				a.checkExpr(spec.Value, scope)
			}
			if prev := scope.Insert(spec.Name, &Symbol{Name: spec.Name, Kind: "const", Decl: spec}); prev != nil {
				a.errorf(spec.Start, "const %q redeclared in this block", spec.Name)
			}
		}
	}
}

// inferCompositeLitType recursively pushes the expected type down into a composite literal.
// In checker.go
func (a *Analyzer) inferCompositeLitType(expr ast.Expr, typeRef ast.TypeRef) {
	if expr == nil || typeRef == nil {
		return
	}

	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return
	}
	if lit.Type == nil {
		lit.Type = typeRef
	}

	switch t := typeRef.(type) {
	case *ast.ArrayType:
		for _, f := range lit.Fields {
			val := f
			if kv, ok := f.(*ast.KeyValueExpr); ok {
				val = kv.Value
			}
			// Recursively infer type for nested array elements
			a.inferCompositeLitType(val, t.Elem)
		}
	case *ast.VectorType:
		for _, f := range lit.Fields {
			val := f
			if kv, ok := f.(*ast.KeyValueExpr); ok {
				val = kv.Value
			}
			// Recursively infer type for nested vector elements
			a.inferCompositeLitType(val, t.Elem)
		}
	case *ast.NamedType:
		// For struct types, we might need to infer nested types
		// but we'd need to look up the struct definition
	}
}

func (a *Analyzer) checkAssignStmt(s *ast.AssignStmt, scope *Scope) {
	a.checkExpr(s.Target, scope)
	if s.Value != nil {
		a.checkExpr(s.Value, scope)
	}
}

func (a *Analyzer) checkForStmt(s *ast.ForStmt, scope *Scope) {
	loopScope := NewScope(scope, ScopeBlock)
	if s.Init != nil {
		a.checkStmt(s.Init, loopScope)
	}
	if s.Cond != nil {
		a.checkExpr(s.Cond, loopScope)
	}
	if s.Post != nil {
		a.checkStmt(s.Post, loopScope)
	}
	a.checkBlock(s.Body, loopScope)
}

func (a *Analyzer) checkExpr(expr ast.Expr, scope *Scope) {
	if expr == nil {
		return
	}
	switch e := expr.(type) {
	case *ast.Ident:
		if scope.Lookup(e.Name) == nil {
			a.errorf(e.NamePos, "undefined: %s", e.Name)
		}
	case *ast.BinaryExpr:
		a.checkExpr(e.Left, scope)
		a.checkExpr(e.Right, scope)
	case *ast.UnaryExpr:
		a.checkExpr(e.X, scope)
	case *ast.CallExpr:
		a.checkExpr(e.Fun, scope)
		for _, arg := range e.Args {
			a.checkExpr(arg, scope)
		}
	case *ast.SelectorExpr:
		a.checkExpr(e.X, scope)
	case *ast.IndexExpr:
		a.checkExpr(e.X, scope)
		a.checkExpr(e.Index, scope)
	case *ast.SliceExpr:
		a.checkExpr(e.X, scope)
		a.checkExpr(e.Low, scope)
		a.checkExpr(e.High, scope)
	case *ast.RangeExpr:
		a.checkExpr(e.Low, scope)
		a.checkExpr(e.High, scope)
	case *ast.AwaitExpr:
		a.checkExpr(e.X, scope)
	case *ast.CompositeLit:
		for _, f := range e.Fields {
			a.checkExpr(f, scope)
		}
	case *ast.KeyValueExpr:
		a.checkExpr(e.Value, scope)
	case *ast.TupleLit:
		for _, el := range e.Elems {
			a.checkExpr(el, scope)
		}
	case *ast.LambdaExpr:
		lamScope := NewScope(scope, ScopeFunc)
		for _, p := range e.Params {
			lamScope.Insert(p.Name, &Symbol{Name: p.Name, Kind: "param", Decl: p, Type: p.Type})
		}
		a.checkBlock(e.Body, lamScope)
	case *ast.ProcessExpr:
		procScope := NewScope(scope, ScopeFunc)
		for _, p := range e.Params {
			procScope.Insert(p.Name, &Symbol{Name: p.Name, Kind: "param", Decl: p, Type: p.Type})
		}
		for _, arg := range e.Args {
			a.checkExpr(arg, scope)
		}
		a.checkBlock(e.Body, procScope)
	case *ast.NewExpr:
		a.checkExpr(e.Init, scope)
	case *ast.NewArrayExpr:
		a.checkExpr(e.Len, scope)
	case *ast.DeleteExpr:
		a.checkExpr(e.X, scope)
	case *ast.CastExpr:
		a.checkExpr(e.X, scope)
	}
}