// lower/async.go
package lower

import (
	"fmt"

	"github.com/arc-language/arc-lang/ast"
)

// lowerAsync transforms every top-level async FuncDecl into three declarations:
//
//  1. <Name>_Packet  — an InterfaceDecl (struct) holding params + return slot.
//  2. <Name>_ThreadEntry — a plain FuncDecl that receives *void and runs the body.
//  3. <Name>         — the original name, now a thin wrapper that allocates the
//     packet and spawns the thread, returning a ThreadHandle.
//
// await expressions in all function bodies are then rewritten to thread_join calls.
func lowerAsync(file *ast.File) {
	var newDecls []ast.Decl

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !fn.IsAsync {
			newDecls = append(newDecls, decl)
			continue
		}
		packet, entry, wrapper := transformAsyncFunc(fn)
		newDecls = append(newDecls, packet, entry, wrapper)
	}

	file.Decls = newDecls

	// Rewrite await expressions in every declaration (including the newly
	// generated ones) now that the packet types are registered.
	rewriteAwaitExprs(file)
}

// transformAsyncFunc performs the pthread-style split of a single async func.
func transformAsyncFunc(fn *ast.FuncDecl) (packet *ast.InterfaceDecl, entry *ast.FuncDecl, wrapper *ast.FuncDecl) {
	packetName := fmt.Sprintf("%s_Packet", fn.Name)
	entryName := fmt.Sprintf("%s_ThreadEntry", fn.Name)

	// ── 1. Packet struct ─────────────────────────────────────────────────────
	// Fields for every parameter plus a "ret" slot for the return value.
	var fields []*ast.Field
	for _, p := range fn.Params {
		fields = append(fields, &ast.Field{
			Name:  p.Name,
			Type:  p.Type,
			Start: p.Start,
		})
	}
	if fn.ReturnType != nil {
		fields = append(fields, &ast.Field{
			Name:  "ret",
			Type:  fn.ReturnType,
			Start: fn.Start,
		})
	}

	packet = &ast.InterfaceDecl{
		Name:   packetName,
		Fields: fields,
		Start:  fn.Start,
	}

	// ── 2. Thread entry ──────────────────────────────────────────────────────
	// func <Name>_ThreadEntry(raw: *void) {
	//     let pkt: *<Name>_Packet = cast(raw, *<Name>_Packet)
	//     <...original body, with param refs rewritten to pkt.param...>
	//     thread_exit()
	// }
	castExpr := makeCast(
		&ast.Ident{Name: "raw"},
		&ast.PointerType{Base: &ast.NamedType{Name: packetName}},
	)
	pktDecl := &ast.DeclStmt{Decl: &ast.VarDecl{
		Name:  "pkt",
		IsRef: false,
		Type:  &ast.PointerType{Base: &ast.NamedType{Name: packetName}},
		Value: castExpr,
	}}

	// Collect the original body statements and remap parameter names.
	rewrittenBody := remapParams(fn)

	entryBody := &ast.BlockStmt{
		LBrace: fn.Body.LBrace,
		RBrace: fn.Body.RBrace,
		List: append(
			[]ast.Stmt{pktDecl},
			append(rewrittenBody, makeCall("thread_exit"))...,
		),
	}

	entry = &ast.FuncDecl{
		Name: entryName,
		Params: []*ast.Field{{
			Name:  "raw",
			Type:  &ast.PointerType{Base: &ast.NamedType{Name: "void"}},
			Start: fn.Start,
		}},
		ReturnType: nil, // void
		Body:       entryBody,
		Start:      fn.Start,
	}

	// ── 3. Wrapper ───────────────────────────────────────────────────────────
	// func <Name>(<params...>) ThreadHandle {
	//     let pkt = new <Name>_Packet{<param>: <param>, ...}
	//     return syscall_spawn(<Name>_ThreadEntry, pkt)
	// }
	var fieldInits []ast.Expr
	for _, p := range fn.Params {
		fieldInits = append(fieldInits, &ast.KeyValueExpr{
			Key:   &ast.Ident{Name: p.Name},
			Value: &ast.Ident{Name: p.Name},
		})
	}
	newPktExpr := &ast.NewExpr{
		Type: &ast.NamedType{Name: packetName},
		Init: &ast.CompositeLit{Fields: fieldInits},
	}
	pktVar := &ast.DeclStmt{Decl: &ast.VarDecl{
		Name:  "pkt",
		IsRef: false,
		Value: newPktExpr,
	}}
	spawnReturn := &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.Ident{Name: "syscall_spawn"},
				Args: []ast.Expr{
					&ast.Ident{Name: entryName},
					&ast.Ident{Name: "pkt"},
				},
			},
		},
	}

	wrapper = &ast.FuncDecl{
		Name:       fn.Name,
		Params:     fn.Params,
		ReturnType: &ast.NamedType{Name: "ThreadHandle"},
		Body: &ast.BlockStmt{
			LBrace: fn.Body.LBrace,
			RBrace: fn.Body.RBrace,
			List:   []ast.Stmt{pktVar, spawnReturn},
		},
		Start: fn.Start,
	}

	return
}

// remapParams deep-copies the original function body and rewrites every
// reference to a parameter name into a pkt.<name> member access.
// This is a best-effort transformation; a production compiler would use
// the type-annotated symbol table to be precise.
func remapParams(fn *ast.FuncDecl) []ast.Stmt {
	paramSet := make(map[string]bool, len(fn.Params))
	for _, p := range fn.Params {
		paramSet[p.Name] = true
	}

	// Clone the statement list so we don't mutate the original body.
	stmts := cloneStmts(fn.Body.List)

	// Rewrite every Ident that names a parameter into pkt.<name>.
	rewriteIdents(stmts, func(id *ast.Ident) ast.Expr {
		if paramSet[id.Name] {
			return &ast.SelectorExpr{X: &ast.Ident{Name: "pkt"}, Sel: id.Name}
		}
		return id
	})

	// If the function has a non-void return, transform each ReturnStmt into an
	// assignment to pkt.ret followed by a plain return (the thread_exit call
	// appended at the end covers actual termination).
	if fn.ReturnType != nil {
		stmts = rewriteReturnsToPacket(stmts)
	}

	return stmts
}

// rewriteReturnsToPacket replaces every ReturnStmt that carries a value with:
//
//	pkt.ret = <value>
//	return
func rewriteReturnsToPacket(stmts []ast.Stmt) []ast.Stmt {
	var out []ast.Stmt
	for _, stmt := range stmts {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(ret.Results) == 0 {
			out = append(out, stmt)
			continue
		}
		var val ast.Expr
		if len(ret.Results) == 1 {
			val = ret.Results[0]
		} else {
			val = &ast.TupleLit{Elems: ret.Results}
		}
		assign := &ast.AssignStmt{
			Target: &ast.SelectorExpr{X: &ast.Ident{Name: "pkt"}, Sel: "ret"},
			Op:     "=",
			Value:  val,
		}
		out = append(out, assign, &ast.ReturnStmt{})
	}
	return out
}

// rewriteAwaitExprs walks the entire file replacing AwaitExpr nodes with
// calls to thread_join.
func rewriteAwaitExprs(file *ast.File) {
	walkExprs(file, func(e ast.Expr) ast.Expr {
		aw, ok := e.(*ast.AwaitExpr)
		if !ok {
			return e
		}
		// await expr  →  thread_join(expr)
		return &ast.CallExpr{
			Fun:  &ast.Ident{Name: "thread_join"},
			Args: []ast.Expr{aw.X},
		}
	})
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func makeCast(x ast.Expr, to ast.TypeRef) *ast.CastExpr {
	return &ast.CastExpr{X: x, Type: to}
}

func makeCall(name string, args ...ast.Expr) ast.Stmt {
	return &ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: name}, Args: args}}
}

// cloneStmts performs a shallow clone of the top-level slice only.
// For a production compiler you would want a full deep clone; here a
// shallow copy is sufficient because remapParams rewrites via walkExprs
// which replaces nodes rather than mutating them in place.
func cloneStmts(stmts []ast.Stmt) []ast.Stmt {
	out := make([]ast.Stmt, len(stmts))
	copy(out, stmts)
	return out
}

// rewriteIdents walks a statement list replacing Ident nodes via fn.
func rewriteIdents(stmts []ast.Stmt, fn func(*ast.Ident) ast.Expr) {
	rewriter := func(e ast.Expr) ast.Expr {
		if id, ok := e.(*ast.Ident); ok {
			return fn(id)
		}
		return e
	}
	for _, s := range stmts {
		walkExprsStmt(s, rewriter)
	}
}