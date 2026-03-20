// frontend/translator.go
package frontend

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/arc-language/arc-lang/ast"
	"github.com/arc-language/arc-lang/parser"
)

// Translate converts the raw ANTLR CST into a clean Arc AST.
func Translate(root parser.ICompilationUnitContext) *ast.File {
	t := &translator{}
	return t.file(root)
}

type translator struct{}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func (t *translator) pos(token antlr.Token) ast.Position {
	if token == nil {
		return ast.Position{}
	}
	return ast.Position{Line: token.GetLine(), Column: token.GetColumn()}
}

// stripQuotes removes the surrounding double-quotes from a STRING_LIT value.
func stripQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// ─── File ────────────────────────────────────────────────────────────────────

func (t *translator) file(ctx parser.ICompilationUnitContext) *ast.File {
	f := &ast.File{Start: t.pos(ctx.GetStart())}

	if ns := ctx.NamespaceDecl(); ns != nil {
		parts := make([]string, 0, len(ns.AllIDENTIFIER()))
		for _, id := range ns.AllIDENTIFIER() {
			parts = append(parts, id.GetText())
		}
		f.Namespace = strings.Join(parts, ".")
	}

	for _, tld := range ctx.AllTopLevelDecl() {
		nodes := t.topLevelDecl(tld)
		for _, n := range nodes {
			switch v := n.(type) {
			case *ast.ImportDecl:
				f.Imports = append(f.Imports, v)
			default:
				f.Decls = append(f.Decls, v)
			}
		}
	}
	return f
}

// topLevelDecl returns a slice because a grouped import "import ( ... )"
// produces multiple ImportDecl nodes from one TopLevelDeclContext.
func (t *translator) topLevelDecl(ctx parser.ITopLevelDeclContext) []ast.Decl {
	// Collect leading attributes (used by interface and any future attributed decl).
	var attrs []*ast.Attribute
	for _, a := range ctx.AllAttribute() {
		attrs = append(attrs, t.attribute(a))
	}

	if ctx.ImportDecl() != nil {
		return t.importDecl(ctx.ImportDecl())
	}
	if ctx.ConstDecl() != nil {
		return []ast.Decl{t.constDecl(ctx.ConstDecl())}
	}
	if ctx.TopLevelVarDecl() != nil {
		return []ast.Decl{t.topLevelVarDecl(ctx.TopLevelVarDecl())}
	}
	if ctx.TopLevelLetDecl() != nil {
		return []ast.Decl{t.topLevelLetDecl(ctx.TopLevelLetDecl())}
	}
	if ctx.FuncDecl() != nil {
		return []ast.Decl{t.funcDecl(ctx.FuncDecl())}
	}
	if ctx.DeinitDecl() != nil {
		return []ast.Decl{t.deinitDecl(ctx.DeinitDecl())}
	}
	if ctx.InterfaceDecl() != nil {
		d := t.interfaceDecl(ctx.InterfaceDecl())
		d.Attrs = attrs
		return []ast.Decl{d}
	}
	if ctx.EnumDecl() != nil {
		return []ast.Decl{t.enumDecl(ctx.EnumDecl())}
	}
	if ctx.TypeAliasDecl() != nil {
		return []ast.Decl{t.typeAliasDecl(ctx.TypeAliasDecl())}
	}
	if ctx.ExternDecl() != nil {
		return []ast.Decl{t.externDecl(ctx.ExternDecl())}
	}
	return nil
}

// ─── Attributes ──────────────────────────────────────────────────────────────

func (t *translator) attribute(ctx parser.IAttributeContext) *ast.Attribute {
	a := &ast.Attribute{
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.Expression() != nil {
		a.Arg = t.expr(ctx.Expression())
	}
	return a
}

// ─── Imports ─────────────────────────────────────────────────────────────────

// importDecl may produce multiple ImportDecl nodes (grouped import block).
func (t *translator) importDecl(ctx parser.IImportDeclContext) []ast.Decl {
	var result []ast.Decl
	for _, spec := range ctx.AllImportSpec() {
		path := stripQuotes(spec.STRING_LIT().GetText())
		alias := ""
		if spec.ImportAlias() != nil {
			alias = spec.ImportAlias().GetText()
		}
		result = append(result, &ast.ImportDecl{
			Alias: alias,
			Path:  path,
			Start: t.pos(spec.GetStart()),
		})
	}
	return result
}

// ─── Constants ───────────────────────────────────────────────────────────────

func (t *translator) constDecl(ctx parser.IConstDeclContext) *ast.ConstDecl {
	d := &ast.ConstDecl{Start: t.pos(ctx.GetStart())}
	for _, spec := range ctx.AllConstSpec() {
		cs := &ast.ConstSpec{
			Name:  spec.IDENTIFIER().GetText(),
			Value: t.expr(spec.Expression()),
			Start: t.pos(spec.GetStart()),
		}
		if spec.TypeRef() != nil {
			cs.Type = t.typeRef(spec.TypeRef())
		}
		d.Specs = append(d.Specs, cs)
	}
	return d
}

// ─── Variables ───────────────────────────────────────────────────────────────

func (t *translator) topLevelVarDecl(ctx parser.ITopLevelVarDeclContext) *ast.VarDecl {
	d := &ast.VarDecl{
		Name:  ctx.IDENTIFIER().GetText(),
		IsRef: true,
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.TypeRef() != nil {
		d.Type = t.typeRef(ctx.TypeRef())
	}
	if ctx.NULL() != nil {
		d.IsNull = true
	} else if ctx.Expression() != nil {
		d.Value = t.expr(ctx.Expression())
	}
	return d
}

func (t *translator) topLevelLetDecl(ctx parser.ITopLevelLetDeclContext) *ast.VarDecl {
	d := &ast.VarDecl{
		Name:  ctx.IDENTIFIER().GetText(),
		IsRef: false,
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.TypeRef() != nil {
		d.Type = t.typeRef(ctx.TypeRef())
	}
	if ctx.Expression() != nil {
		d.Value = t.expr(ctx.Expression())
	}
	return d
}

// ─── Functions ───────────────────────────────────────────────────────────────

func (t *translator) funcDecl(ctx parser.IFuncDeclContext) *ast.FuncDecl {
	fn := &ast.FuncDecl{
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}

	for _, mod := range ctx.AllFuncModifier() {
		if mod.ASYNC() != nil {
			fn.IsAsync = true
		}
		if mod.GPU() != nil {
			fn.IsGpu = true
		}
	}

	if gp := ctx.GenericParams(); gp != nil {
		for _, id := range gp.AllIDENTIFIER() {
			fn.GenericParams = append(fn.GenericParams, id.GetText())
		}
	}

	if pl := ctx.ParamList(); pl != nil {
		for _, p := range pl.AllParam() {
			// self param — first param may be a self receiver
			if p.SelfParam() != nil {
				fn.Self = t.selfParam(p.SelfParam())
				continue
			}
			if p.ELLIPSIS() != nil {
				fn.IsVariadic = true
				continue
			}
			fn.Params = append(fn.Params, t.param(p))
		}
	}

	if rt := ctx.ReturnType(); rt != nil {
		fn.ReturnType = t.returnType(rt)
	}

	if ctx.Block() != nil {
		fn.Body = t.blockStmt(ctx.Block())
	}
	return fn
}

func (t *translator) deinitDecl(ctx parser.IDeinitDeclContext) *ast.DeinitDecl {
	return &ast.DeinitDecl{
		Self:  t.selfParam(ctx.SelfParam()),
		Body:  t.blockStmt(ctx.Block()),
		Start: t.pos(ctx.GetStart()),
	}
}

func (t *translator) selfParam(ctx parser.ISelfParamContext) *ast.SelfParam {
	sp := &ast.SelfParam{
		Name:  ctx.IDENTIFIER().GetText(),
		IsMut: ctx.MUT() != nil,
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.ParamType() != nil {
		sp.Type = t.paramType(ctx.ParamType())
	}
	return sp
}

func (t *translator) param(ctx parser.IParamContext) *ast.Field {
	f := &ast.Field{
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.ParamType() != nil {
		f.Type = t.paramType(ctx.ParamType())
	}
	return f
}

func (t *translator) paramType(ctx parser.IParamTypeContext) ast.TypeRef {
	if ctx.AMP() != nil && ctx.MUT() != nil {
		return &ast.MutRefType{
			Base:  t.typeRef(ctx.TypeRef()),
			Start: t.pos(ctx.GetStart()),
		}
	}
	return t.typeRef(ctx.TypeRef())
}

func (t *translator) returnType(ctx parser.IReturnTypeContext) ast.TypeRef {
	if ctx.TupleType() != nil {
		return t.tupleType(ctx.TupleType())
	}
	return t.typeRef(ctx.TypeRef())
}

func (t *translator) tupleType(ctx parser.ITupleTypeContext) *ast.TupleType {
	tt := &ast.TupleType{Start: t.pos(ctx.GetStart())}
	for _, tr := range ctx.AllTypeRef() {
		tt.Elems = append(tt.Elems, t.typeRef(tr))
	}
	return tt
}

// ─── Interfaces ───────────────────────────────────────────────────────────────

func (t *translator) interfaceDecl(ctx parser.IInterfaceDeclContext) *ast.InterfaceDecl {
	d := &ast.InterfaceDecl{
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	if gp := ctx.GenericParams(); gp != nil {
		for _, id := range gp.AllIDENTIFIER() {
			d.GenericParams = append(d.GenericParams, id.GetText())
		}
	}
	for _, f := range ctx.AllInterfaceField() {
		d.Fields = append(d.Fields, &ast.Field{
			Name:  f.IDENTIFIER().GetText(),
			Type:  t.typeRef(f.TypeRef()),
			Start: t.pos(f.GetStart()),
		})
	}
	return d
}

// ─── Enums ───────────────────────────────────────────────────────────────────

func (t *translator) enumDecl(ctx parser.IEnumDeclContext) *ast.EnumDecl {
	d := &ast.EnumDecl{
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.PrimitiveType() != nil {
		d.UnderlyingType = &ast.NamedType{
			Name:    ctx.PrimitiveType().GetText(),
			TypePos: t.pos(ctx.PrimitiveType().GetStart()),
		}
	}
	for _, m := range ctx.AllEnumMember() {
		member := &ast.EnumMember{
			Name:  m.IDENTIFIER().GetText(),
			Start: t.pos(m.GetStart()),
		}
		if m.Expression() != nil {
			member.Value = t.expr(m.Expression())
		}
		d.Members = append(d.Members, member)
	}
	return d
}

// ─── Type Aliases ─────────────────────────────────────────────────────────────

func (t *translator) typeAliasDecl(ctx parser.ITypeAliasDeclContext) *ast.TypeAliasDecl {
	d := &ast.TypeAliasDecl{
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.OPAQUE() != nil {
		d.IsOpaque = true
	} else if ctx.TypeRef() != nil {
		d.Type = t.typeRef(ctx.TypeRef())
	}
	return d
}

// ─── Extern ───────────────────────────────────────────────────────────────────

func (t *translator) externDecl(ctx parser.IExternDeclContext) *ast.ExternDecl {
	d := &ast.ExternDecl{
		Lang:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	for _, m := range ctx.AllExternMember() {
		if em := t.externMember(m); em != nil {
			d.Members = append(d.Members, em)
		}
	}
	return d
}

func (t *translator) externMember(ctx parser.IExternMemberContext) ast.ExternMember {
	if ctx.ExternFuncDecl() != nil {
		return t.externFunc(ctx.ExternFuncDecl())
	}
	if ctx.ExternTypeAlias() != nil {
		return t.externTypeAlias(ctx.ExternTypeAlias())
	}
	if ctx.ExternNamespace() != nil {
		return t.externNamespace(ctx.ExternNamespace())
	}
	if ctx.ExternClass() != nil {
		return t.externClass(ctx.ExternClass())
	}
	return nil
}

func (t *translator) externFunc(ctx parser.IExternFuncDeclContext) *ast.ExternFunc {
	ef := &ast.ExternFunc{
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.CallingConvention() != nil {
		ef.Convention = ctx.CallingConvention().GetText()
	}
	if ctx.ExternSymbol() != nil {
		ef.Symbol = stripQuotes(ctx.ExternSymbol().STRING_LIT().GetText())
	}
	if ctx.ExternParamList() != nil {
		pl := ctx.ExternParamList()
		for _, p := range pl.AllExternParam() {
			ef.Params = append(ef.Params, t.externType(p.ExternType()))
		}
		if pl.ELLIPSIS() != nil {
			ef.IsVariadic = true
		}
	}
	if ctx.ExternReturnType() != nil {
		ef.Return = t.externType(ctx.ExternReturnType().ExternType())
	}
	return ef
}

func (t *translator) externTypeAlias(ctx parser.IExternTypeAliasContext) *ast.ExternTypeAlias {
	ea := &ast.ExternTypeAlias{
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	fpt := ctx.ExternFunctionPtrType()
	ft := &ast.FuncType{Start: t.pos(fpt.GetStart())}
	if fpt.ExternParamList() != nil {
		for _, p := range fpt.ExternParamList().AllExternParam() {
			ft.Params = append(ft.Params, t.externType(p.ExternType()))
		}
	}
	if fpt.ExternReturnType() != nil {
		ft.Results = []ast.TypeRef{t.externType(fpt.ExternReturnType().ExternType())}
	}
	ea.Type = ft
	return ea
}

func (t *translator) externNamespace(ctx parser.IExternNamespaceContext) *ast.ExternNamespace {
	parts := make([]string, 0, len(ctx.AllIDENTIFIER()))
	for _, id := range ctx.AllIDENTIFIER() {
		parts = append(parts, id.GetText())
	}
	en := &ast.ExternNamespace{
		Name:  strings.Join(parts, "."),
		Start: t.pos(ctx.GetStart()),
	}
	for _, m := range ctx.AllExternMember() {
		if em := t.externMember(m); em != nil {
			en.Members = append(en.Members, em)
		}
	}
	return en
}

func (t *translator) externClass(ctx parser.IExternClassContext) *ast.ExternClass {
	ec := &ast.ExternClass{
		IsAbstract: ctx.ABSTRACT() != nil,
		Name:       ctx.IDENTIFIER().GetText(),
		Start:      t.pos(ctx.GetStart()),
	}
	if ctx.ExternSymbol() != nil {
		ec.Symbol = stripQuotes(ctx.ExternSymbol().STRING_LIT().GetText())
	}
	for _, m := range ctx.AllExternClassMember() {
		if em := t.externClassMember(m); em != nil {
			ec.Methods = append(ec.Methods, em)
		}
	}
	return ec
}

func (t *translator) externClassMember(ctx parser.IExternClassMemberContext) *ast.ExternMethod {
	if ctx.ExternVirtualMethod() != nil {
		return t.externVirtualMethod(ctx.ExternVirtualMethod())
	}
	if ctx.ExternStaticMethod() != nil {
		return t.externStaticMethod(ctx.ExternStaticMethod())
	}
	if ctx.ExternConstructor() != nil {
		return t.externConstructor(ctx.ExternConstructor())
	}
	if ctx.ExternDestructor() != nil {
		return t.externDestructor(ctx.ExternDestructor())
	}
	return nil
}

func (t *translator) externVirtualMethod(ctx parser.IExternVirtualMethodContext) *ast.ExternMethod {
	m := &ast.ExternMethod{
		Kind:  ast.ExternVirtual,
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.CallingConvention() != nil {
		m.Convention = ctx.CallingConvention().GetText()
	}
	if ctx.ExternMethodParamList() != nil {
		pl := ctx.ExternMethodParamList()
		// self is first
		if pl.ExternMethodParam() != nil {
			m.Params = append(m.Params, t.externType(pl.ExternMethodParam().ExternType()))
		}
		for _, p := range pl.AllExternParam() {
			m.Params = append(m.Params, t.externType(p.ExternType()))
		}
		if pl.ELLIPSIS() != nil {
			m.IsVariadic = true
		}
	}
	if ctx.ExternReturnType() != nil {
		rt := ctx.ExternReturnType()
		m.IsConst = rt.CONST() != nil
		m.Return = t.externType(rt.ExternType())
	}
	return m
}

func (t *translator) externStaticMethod(ctx parser.IExternStaticMethodContext) *ast.ExternMethod {
	m := &ast.ExternMethod{
		Kind:  ast.ExternStatic,
		Name:  ctx.IDENTIFIER().GetText(),
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.ExternSymbol() != nil {
		m.Symbol = stripQuotes(ctx.ExternSymbol().STRING_LIT().GetText())
	}
	if ctx.ExternParamList() != nil {
		pl := ctx.ExternParamList()
		for _, p := range pl.AllExternParam() {
			m.Params = append(m.Params, t.externType(p.ExternType()))
		}
		if pl.ELLIPSIS() != nil {
			m.IsVariadic = true
		}
	}
	if ctx.ExternReturnType() != nil {
		m.Return = t.externType(ctx.ExternReturnType().ExternType())
	}
	return m
}

func (t *translator) externConstructor(ctx parser.IExternConstructorContext) *ast.ExternMethod {
	m := &ast.ExternMethod{
		Kind:   ast.ExternConstructor,
		Name:   "new",
		Return: t.externType(ctx.ExternType()),
		Start:  t.pos(ctx.GetStart()),
	}
	if ctx.ExternParamList() != nil {
		for _, p := range ctx.ExternParamList().AllExternParam() {
			m.Params = append(m.Params, t.externType(p.ExternType()))
		}
	}
	return m
}

func (t *translator) externDestructor(ctx parser.IExternDestructorContext) *ast.ExternMethod {
	m := &ast.ExternMethod{
		Kind:  ast.ExternDestructor,
		Name:  "delete",
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.ExternMethodParam() != nil {
		m.Params = append(m.Params, t.externType(ctx.ExternMethodParam().ExternType()))
	}
	return m
}

// externType recursively translates the extern type grammar.
// The grammar is left-recursive: **T, *T, *const T, &T, &const T, [N]T, primitives, named.
func (t *translator) externType(ctx parser.IExternTypeContext) ast.TypeRef {
	start := t.pos(ctx.GetStart())

	// **T  (two STAR tokens, inner ExternType)
	stars := ctx.AllSTAR()
	if len(stars) == 2 {
		inner := &ast.PointerType{
			IsConst: ctx.CONST() != nil,
			Base:    t.externType(ctx.ExternType()),
			Start:   start,
		}
		return &ast.PointerType{Base: inner, Start: start}
	}
	// *T or *const T
	if len(stars) == 1 {
		return &ast.PointerType{
			IsConst: ctx.CONST() != nil,
			Base:    t.externType(ctx.ExternType()),
			Start:   start,
		}
	}
	// &T or &const T
	if ctx.AMP() != nil {
		return &ast.RefType{
			IsConst: ctx.CONST() != nil,
			Base:    t.externType(ctx.ExternType()),
			Start:   start,
		}
	}
	// [N]T
	if ctx.LBRACKET() != nil {
		return &ast.ArrayType{
			Len:   t.expr(ctx.Expression()),
			Elem:  t.externType(ctx.ExternType()),
			Start: start,
		}
	}
	// Leaf types
	if ctx.PrimitiveType() != nil {
		return &ast.NamedType{Name: ctx.PrimitiveType().GetText(), TypePos: start}
	}
	if ctx.QualifiedName() != nil {
		return &ast.NamedType{Name: ctx.QualifiedName().GetText(), TypePos: start}
	}
	// void, bool, string, byte, char, usize, isize
	return &ast.NamedType{Name: ctx.GetText(), TypePos: start}
}

// ─── Statements ───────────────────────────────────────────────────────────────

func (t *translator) blockStmt(ctx parser.IBlockContext) *ast.BlockStmt {
	b := &ast.BlockStmt{
		LBrace: t.pos(ctx.LBRACE().GetSymbol()),
		RBrace: t.pos(ctx.RBRACE().GetSymbol()),
	}
	for _, s := range ctx.AllStatement() {
		if st := t.stmt(s); st != nil {
			b.List = append(b.List, st)
		}
	}
	return b
}

func (t *translator) stmt(ctx parser.IStatementContext) ast.Stmt {
	if ctx.LetStatement() != nil {
		return t.letStmt(ctx.LetStatement())
	}
	if ctx.VarStatement() != nil {
		return t.varStmt(ctx.VarStatement())
	}
	if ctx.ConstDecl() != nil {
		return &ast.DeclStmt{Decl: t.constDecl(ctx.ConstDecl())}
	}
	if ctx.ReturnStatement() != nil {
		return t.returnStmt(ctx.ReturnStatement())
	}
	if ctx.BreakStatement() != nil {
		return &ast.BreakStmt{Start: t.pos(ctx.GetStart())}
	}
	if ctx.ContinueStatement() != nil {
		return &ast.ContinueStmt{Start: t.pos(ctx.GetStart())}
	}
	if ctx.DeferStatement() != nil {
		return t.deferStmt(ctx.DeferStatement())
	}
	if ctx.IfStatement() != nil {
		return t.ifStmt(ctx.IfStatement())
	}
	if ctx.ForStatement() != nil {
		return t.forStmt(ctx.ForStatement())
	}
	if ctx.SwitchStatement() != nil {
		return t.switchStmt(ctx.SwitchStatement())
	}
	if ctx.AssignmentStatement() != nil {
		return t.assignStmt(ctx.AssignmentStatement())
	}
	if ctx.ExpressionStatement() != nil {
		return &ast.ExprStmt{X: t.expr(ctx.ExpressionStatement().Expression())}
	}
	return nil
}

func (t *translator) letStmt(ctx parser.ILetStatementContext) ast.Stmt {
	start := t.pos(ctx.GetStart())

	// Destructuring: let (a, b) = expr
	if ctx.LPAREN() != nil {
		ids := ctx.AllIDENTIFIER()
		// Build a multi-bind as a DeclStmt holding a VarDecl whose Name encodes
		// the tuple. The lower/checker can unpack it from Value being a TupleLit.
		// For now we emit one DeclStmt per name — the checker will handle it.
		// A cleaner model would be a dedicated DestructureStmt; add that later if needed.
		decl := &ast.VarDecl{
			Name:  strings.Join(identTexts(ids), ","),
			IsRef: false,
			Value: t.expr(ctx.Expression()),
			Start: start,
		}
		return &ast.DeclStmt{Decl: decl}
	}

	decl := &ast.VarDecl{
		Name:  ctx.IDENTIFIER(0).GetText(),
		IsRef: false,
		Start: start,
	}
	if ctx.TypeRef() != nil {
		decl.Type = t.typeRef(ctx.TypeRef())
	}
	if ctx.Expression() != nil {
		decl.Value = t.expr(ctx.Expression())
	}
	return &ast.DeclStmt{Decl: decl}
}

func (t *translator) varStmt(ctx parser.IVarStatementContext) ast.Stmt {
	decl := &ast.VarDecl{
		Name:  ctx.IDENTIFIER().GetText(),
		IsRef: true,
		Start: t.pos(ctx.GetStart()),
	}
	if ctx.TypeRef() != nil {
		decl.Type = t.typeRef(ctx.TypeRef())
	}
	if ctx.NULL() != nil {
		decl.IsNull = true
	} else if ctx.Expression() != nil {
		decl.Value = t.expr(ctx.Expression())
	}
	return &ast.DeclStmt{Decl: decl}
}

func (t *translator) returnStmt(ctx parser.IReturnStatementContext) *ast.ReturnStmt {
	r := &ast.ReturnStmt{Start: t.pos(ctx.GetStart())}
	for _, e := range ctx.AllExpression() {
		r.Results = append(r.Results, t.expr(e))
	}
	return r
}

func (t *translator) deferStmt(ctx parser.IDeferStatementContext) *ast.DeferStmt {
	return &ast.DeferStmt{
		Call:  t.expr(ctx.Expression()),
		Start: t.pos(ctx.GetStart()),
	}
}

func (t *translator) ifStmt(ctx parser.IIfStatementContext) *ast.IfStmt {
	// Grammar: IF expr block (ELSE IF expr block)* (ELSE block)?
	// AllIF, AllExpression, AllBlock, AllELSE are parallel slices.
	exprs := ctx.AllExpression()
	blocks := ctx.AllBlock()

	// Build the chain from the innermost else-if outward.
	// Index 0 = first if, index 1+ = else-if branches.
	var build func(i int) *ast.IfStmt
	build = func(i int) *ast.IfStmt {
		node := &ast.IfStmt{
			Cond:  t.expr(exprs[i]),
			Body:  t.blockStmt(blocks[i]),
			Start: t.pos(ctx.IF(i).GetSymbol()),
		}
		if i+1 < len(exprs) {
			// There is another else-if.
			node.Else = build(i + 1)
		} else if len(blocks) > len(exprs) {
			// There is a trailing else block (no extra expression).
			node.Else = t.blockStmt(blocks[len(blocks)-1])
		}
		return node
	}
	return build(0)
}

func (t *translator) forStmt(ctx parser.IForStatementContext) ast.Stmt {
	start := t.pos(ctx.GetStart())
	hdr := ctx.ForHeader()
	body := t.blockStmt(ctx.Block())

	// for-in: ForIterator is set
	if hdr.ForIterator() != nil {
		fi := hdr.ForIterator()
		ids := fi.AllIDENTIFIER()
		s := &ast.ForInStmt{
			Key:   ids[0].GetText(),
			Iter:  t.expr(fi.Expression()),
			Body:  body,
			Start: start,
		}
		if len(ids) > 1 {
			s.Value = ids[1].GetText()
		}
		return s
	}

	// C-style: two SEMIs present
	if len(hdr.AllSEMI()) == 2 {
		s := &ast.ForStmt{Body: body, Start: start}
		if hdr.ForInit() != nil {
			s.Init = t.forInit(hdr.ForInit())
		}
		if hdr.Expression() != nil {
			s.Cond = t.expr(hdr.Expression())
		}
		if hdr.ForPost() != nil {
			s.Post = t.forPost(hdr.ForPost())
		}
		return s
	}

	// while-style: single expression
	if hdr.Expression() != nil {
		return &ast.ForStmt{
			Cond:  t.expr(hdr.Expression()),
			Body:  body,
			Start: start,
		}
	}

	// infinite
	return &ast.ForStmt{Body: body, Start: start}
}

func (t *translator) forInit(ctx parser.IForInitContext) ast.Stmt {
	if ctx.LET() != nil {
		decl := &ast.VarDecl{
			Name:  ctx.IDENTIFIER().GetText(),
			IsRef: false,
			Value: t.expr(ctx.Expression()),
			Start: t.pos(ctx.GetStart()),
		}
		if ctx.TypeRef() != nil {
			decl.Type = t.typeRef(ctx.TypeRef())
		}
		return &ast.DeclStmt{Decl: decl}
	}
	// expression-as-init (rare but grammar allows it)
	return &ast.ExprStmt{X: t.expr(ctx.Expression())}
}

func (t *translator) forPost(ctx parser.IForPostContext) ast.Stmt {
	start := t.pos(ctx.GetStart())
	if ctx.INC() != nil {
		return &ast.AssignStmt{
			Target: t.expr(ctx.Expression()),
			Op:     "++",
			Start:  start,
		}
	}
	if ctx.DEC() != nil {
		return &ast.AssignStmt{
			Target: t.expr(ctx.Expression()),
			Op:     "--",
			Start:  start,
		}
	}
	if ctx.AssignmentTarget() != nil {
		return &ast.AssignStmt{
			Target: t.assignTarget(ctx.AssignmentTarget()),
			Op:     ctx.AssignOp().GetText(),
			Value:  t.expr(ctx.Expression()),
			Start:  start,
		}
	}
	return &ast.ExprStmt{X: t.expr(ctx.Expression())}
}

func (t *translator) switchStmt(ctx parser.ISwitchStatementContext) *ast.SwitchStmt {
	s := &ast.SwitchStmt{
		Tag:   t.expr(ctx.Expression()),
		Start: t.pos(ctx.GetStart()),
		End_:  t.pos(ctx.RBRACE().GetSymbol()),
	}
	for _, c := range ctx.AllSwitchCase() {
		sc := &ast.SwitchCase{Start: t.pos(c.GetStart())}
		for _, e := range c.ExpressionList().AllExpression() {
			sc.Values = append(sc.Values, t.expr(e))
		}
		for _, st := range c.AllStatement() {
			if stmt := t.stmt(st); stmt != nil {
				sc.Body = append(sc.Body, stmt)
			}
		}
		s.Cases = append(s.Cases, sc)
	}
	if ctx.SwitchDefault() != nil {
		def := ctx.SwitchDefault()
		for _, st := range def.AllStatement() {
			if stmt := t.stmt(st); stmt != nil {
				s.Default = append(s.Default, stmt)
			}
		}
	}
	return s
}

func (t *translator) assignStmt(ctx parser.IAssignmentStatementContext) *ast.AssignStmt {
	s := &ast.AssignStmt{Start: t.pos(ctx.GetStart())}

	if ctx.INC() != nil {
		s.Target = t.assignTarget(ctx.AssignmentTarget())
		s.Op = "++"
		return s
	}
	if ctx.DEC() != nil {
		s.Target = t.assignTarget(ctx.AssignmentTarget())
		s.Op = "--"
		return s
	}
	s.Target = t.assignTarget(ctx.AssignmentTarget())
	s.Op = ctx.AssignOp().GetText()
	s.Value = t.expr(ctx.Expression())
	return s
}

func (t *translator) assignTarget(ctx parser.IAssignmentTargetContext) ast.Expr {
	exprs := ctx.AllExpression()
	if ctx.DOT() != nil {
		// a.b
		return &ast.SelectorExpr{
			X:   t.expr(exprs[0]),
			Sel: ctx.IDENTIFIER().GetText(),
			Dot: t.pos(ctx.DOT().GetSymbol()),
		}
	}
	if ctx.LBRACKET() != nil {
		// a[i]
		return &ast.IndexExpr{
			X:      t.expr(exprs[0]),
			Index:  t.expr(exprs[1]),
			LBrack: t.pos(ctx.LBRACKET().GetSymbol()),
			RBrack: t.pos(ctx.RBRACKET().GetSymbol()),
		}
	}
	// bare identifier
	return &ast.Ident{
		Name:    ctx.IDENTIFIER().GetText(),
		NamePos: t.pos(ctx.GetStart()),
	}
}

// ─── Expressions ─────────────────────────────────────────────────────────────

func (t *translator) expr(ctx parser.IExpressionContext) ast.Expr {
	switch e := ctx.(type) {

	// ── Primary wrapper ──
	case *parser.PrimaryExprContext:
		return t.primary(e.Primary())

	// ── Postfix ──
	case *parser.MemberAccessContext:
		return &ast.SelectorExpr{
			X:   t.expr(e.Expression()),
			Sel: e.IDENTIFIER().GetText(),
			Dot: t.pos(e.DOT().GetSymbol()),
		}
	case *parser.IndexExprContext:
		exprs := e.AllExpression()
		return &ast.IndexExpr{
			X:      t.expr(exprs[0]),
			Index:  t.expr(exprs[1]),
			LBrack: t.pos(e.LBRACKET().GetSymbol()),
			RBrack: t.pos(e.RBRACKET().GetSymbol()),
		}
	case *parser.SliceExprContext:
		exprs := e.AllExpression()
		return &ast.SliceExpr{
			X:      t.expr(exprs[0]),
			Low:    t.expr(exprs[1]),
			High:   t.expr(exprs[2]),
			LBrack: t.pos(e.LBRACKET().GetSymbol()),
			RBrack: t.pos(e.RBRACKET().GetSymbol()),
		}
	case *parser.CallExprContext:
		call := &ast.CallExpr{
			Fun:    t.expr(e.Expression()),
			Start:  t.pos(e.GetStart()),
			EndPos: t.pos(e.GetStop()),
		}
		if e.ArgumentList() != nil {
			for _, arg := range e.ArgumentList().AllArgument() {
				call.Args = append(call.Args, t.expr(arg.Expression()))
			}
		}
		return call
	case *parser.PostIncrementContext:
		return &ast.UnaryExpr{Op: "++", X: t.expr(e.Expression()), Start: t.pos(e.GetStart())}
	case *parser.PostDecrementContext:
		return &ast.UnaryExpr{Op: "--", X: t.expr(e.Expression()), Start: t.pos(e.GetStart())}

	// ── Unary prefix ──
	case *parser.UnaryMinusContext:
		return &ast.UnaryExpr{Op: "-", X: t.expr(e.Expression()), Start: t.pos(e.GetStart())}
	case *parser.LogicalNotContext:
		return &ast.UnaryExpr{Op: "!", X: t.expr(e.Expression()), Start: t.pos(e.GetStart())}
	case *parser.BitwiseNotContext:
		return &ast.UnaryExpr{Op: "~", X: t.expr(e.Expression()), Start: t.pos(e.GetStart())}
	case *parser.AddressOfContext:
		// &x inside arc code is only valid as the argument to memptr or &mut params.
		// Represent it as a unary & for the checker to validate context.
		return &ast.UnaryExpr{Op: "&", X: t.expr(e.Expression()), Start: t.pos(e.GetStart())}
	case *parser.AwaitExprContext:
		return &ast.AwaitExpr{X: t.expr(e.Expression()), Start: t.pos(e.GetStart())}

	// ── Binary ──
	case *parser.MulExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: e.GetOp().GetText(), Right: t.expr(exprs[1]), OpPos: t.pos(e.GetOp())}
	case *parser.AddExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: e.GetOp().GetText(), Right: t.expr(exprs[1]), OpPos: t.pos(e.GetOp())}
	case *parser.ShiftExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: e.GetOp().GetText(), Right: t.expr(exprs[1]), OpPos: t.pos(e.GetOp())}
	case *parser.RelationalExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: e.GetOp().GetText(), Right: t.expr(exprs[1]), OpPos: t.pos(e.GetOp())}
	case *parser.EqualityExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: e.GetOp().GetText(), Right: t.expr(exprs[1]), OpPos: t.pos(e.GetOp())}
	case *parser.BitwiseAndExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: "&", Right: t.expr(exprs[1])}
	case *parser.BitwiseXorExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: "^", Right: t.expr(exprs[1])}
	case *parser.BitwiseOrExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: "|", Right: t.expr(exprs[1])}
	case *parser.LogicalAndExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: "&&", Right: t.expr(exprs[1])}
	case *parser.LogicalOrExprContext:
		exprs := e.AllExpression()
		return &ast.BinaryExpr{Left: t.expr(exprs[0]), Op: "||", Right: t.expr(exprs[1])}
	case *parser.RangeExprContext:
		exprs := e.AllExpression()
		return &ast.RangeExpr{
			Low:    t.expr(exprs[0]),
			High:   t.expr(exprs[1]),
			DotDot: t.pos(e.RANGE().GetSymbol()),
		}
	}
	return nil
}

func (t *translator) primary(ctx parser.IPrimaryContext) ast.Expr {
	switch p := ctx.(type) {

	// ── Literals ──
	case *parser.IntLiteralContext:
		return &ast.BasicLit{Kind: "INT", Value: p.INT_LIT().GetText(), LitPos: t.pos(p.GetStart())}
	case *parser.HexLiteralContext:
		return &ast.BasicLit{Kind: "HEX", Value: p.HEX_LIT().GetText(), LitPos: t.pos(p.GetStart())}
	case *parser.FloatLiteralContext:
		return &ast.BasicLit{Kind: "FLOAT", Value: p.FLOAT_LIT().GetText(), LitPos: t.pos(p.GetStart())}
	case *parser.StringLiteralContext:
		return &ast.BasicLit{Kind: "STRING", Value: p.STRING_LIT().GetText(), LitPos: t.pos(p.GetStart())}
	case *parser.CharLiteralContext:
		return &ast.BasicLit{Kind: "CHAR", Value: p.CHAR_LIT().GetText(), LitPos: t.pos(p.GetStart())}
	case *parser.TrueLiteralContext:
		return &ast.BasicLit{Kind: "BOOL", Value: "true", LitPos: t.pos(p.GetStart())}
	case *parser.FalseLiteralContext:
		return &ast.BasicLit{Kind: "BOOL", Value: "false", LitPos: t.pos(p.GetStart())}
	case *parser.NullLiteralContext:
		return &ast.BasicLit{Kind: "NULL", Value: "null", LitPos: t.pos(p.GetStart())}

	// ── Names ──
	case *parser.IdentExprContext:
		return &ast.Ident{Name: p.IDENTIFIER().GetText(), NamePos: t.pos(p.GetStart())}
	case *parser.QualifiedExprContext:
		return &ast.Ident{Name: p.QualifiedName().GetText(), NamePos: t.pos(p.GetStart())}

	// ── Type in expression position (cast target) ──
	case *parser.PrimitiveTypeExprContext:
		return &ast.Ident{Name: p.PrimitiveType().GetText(), NamePos: t.pos(p.GetStart())}

	// ── Parenthesised ──
	case *parser.ParenExprContext:
		return t.expr(p.Expression())

	// ── Tuple literal ──
	case *parser.TupleLiteralContext:
		tl := &ast.TupleLit{
			LParen: t.pos(p.LPAREN().GetSymbol()),
			RParen: t.pos(p.RPAREN().GetSymbol()),
		}
		for _, e := range p.AllExpression() {
			tl.Elems = append(tl.Elems, t.expr(e))
		}
		return tl

	// ── new ──
	case *parser.NewExprContext:
		return &ast.NewExpr{
			Type:  t.typeRef(p.TypeRef()),
			Init:  t.initializerBlock(p.InitializerBlock()),
			Start: t.pos(p.GetStart()),
		}
	case *parser.NewArrayExprContext:
		return &ast.NewArrayExpr{
			Len:   t.expr(p.Expression()),
			Elem:  t.typeRef(p.TypeRef()),
			Start: t.pos(p.GetStart()),
		}

	// ── delete ──
	case *parser.DeleteExprContext:
		return &ast.DeleteExpr{
			X:     t.expr(p.Expression()),
			Start: t.pos(p.GetStart()),
			End_:  t.pos(p.GetStop()),
		}

	// ── Bare initializer: {...} ──
	case *parser.BareInitExprContext:
		return t.initializerBlock(p.InitializerBlock())

	// ── Typed initializer: Point{x:1} or Box[int32]{...} ──
	case *parser.TypedInitExprContext:
		// FIX: Check if it's a QualifiedName OR a simple IDENTIFIER
		name := ""
		if p.QualifiedName() != nil {
			name = p.QualifiedName().GetText()
		} else {
			name = p.IDENTIFIER().GetText()
		}

		var gargs []ast.TypeRef
		if p.GenericArgs() != nil {
			for _, tr := range p.GenericArgs().AllTypeRef() {
				gargs = append(gargs, t.typeRef(tr))
			}
		}
		typeRef := &ast.NamedType{Name: name, GenericArgs: gargs, TypePos: t.pos(p.GetStart())}
		init := t.initializerBlock(p.InitializerBlock())
		init.Type = typeRef
		return init

	// ── vector[T]{...} ──
	case *parser.VectorLiteralContext:
		elem := t.typeRef(p.TypeRef())
		init := t.initializerBlock(p.InitializerBlock())
		init.Type = &ast.VectorType{Elem: elem, Start: t.pos(p.GetStart())}
		return init

	// ── map[K]V{...} ──
	case *parser.MapLiteralContext:
		typeRefs := p.AllTypeRef()
		init := t.initializerBlock(p.InitializerBlock())
		init.Type = &ast.MapType{
			Key:   t.typeRef(typeRefs[0]),
			Value: t.typeRef(typeRefs[1]),
			Start: t.pos(p.GetStart()),
		}
		return init

	// ── Lambda ──
	case *parser.LambdaExprContext:
		lam := &ast.LambdaExpr{
			IsAsync: p.ASYNC() != nil,
			Body:    t.blockStmt(p.Block()),
			Start:   t.pos(p.GetStart()),
		}
		if p.LambdaParamList() != nil {
			for _, lp := range p.LambdaParamList().AllLambdaParam() {
				lam.Params = append(lam.Params, &ast.Field{
					Name:  lp.IDENTIFIER().GetText(),
					Type:  t.typeRef(lp.TypeRef()),
					Start: t.pos(lp.GetStart()),
				})
			}
		}
		return lam

	// ── process func(...){...}(args) ──
	case *parser.ProcessExprContext:
		pe := &ast.ProcessExpr{
			Body:  t.blockStmt(p.Block()),
			Start: t.pos(p.GetStart()),
		}
		if p.ParamList() != nil {
			for _, param := range p.ParamList().AllParam() {
				if param.ELLIPSIS() == nil && param.SelfParam() == nil {
					pe.Params = append(pe.Params, t.param(param))
				}
			}
		}
		if p.ArgumentList() != nil {
			for _, arg := range p.ArgumentList().AllArgument() {
				pe.Args = append(pe.Args, t.expr(arg.Expression()))
			}
		}
		return pe
	}
	return nil
}

func (t *translator) initializerBlock(ctx parser.IInitializerBlockContext) *ast.CompositeLit {
	lit := &ast.CompositeLit{
		LBrace: t.pos(ctx.LBRACE().GetSymbol()),
		RBrace: t.pos(ctx.RBRACE().GetSymbol()),
	}
	// Field initializers: { key: val }
	if len(ctx.AllFieldInit()) > 0 {
		for _, fi := range ctx.AllFieldInit() {
			lit.Fields = append(lit.Fields, &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: fi.IDENTIFIER().GetText(), NamePos: t.pos(fi.GetStart())},
				Value: t.expr(fi.Expression()),
				Colon: t.pos(fi.COLON().GetSymbol()),
			})
		}
		return lit
	}
	// Map entries: { k: v }
	if len(ctx.AllMapEntry()) > 0 {
		for _, me := range ctx.AllMapEntry() {
			exprs := me.AllExpression()
			lit.Fields = append(lit.Fields, &ast.KeyValueExpr{
				Key:   t.expr(exprs[0]),
				Value: t.expr(exprs[1]),
				Colon: t.pos(me.COLON().GetSymbol()),
			})
		}
		return lit
	}
	// Plain values: { 1, 2, 3 }
	for _, e := range ctx.AllExpression() {
		lit.Fields = append(lit.Fields, t.expr(e))
	}
	return lit
}

// ─── Types ────────────────────────────────────────────────────────────────────

func (t *translator) typeRef(ctx parser.ITypeRefContext) ast.TypeRef {
	if ctx == nil {
		return nil
	}
	if ctx.FunctionType() != nil {
		return t.functionType(ctx.FunctionType())
	}
	return t.baseType(ctx.BaseType())
}

func (t *translator) functionType(ctx parser.IFunctionTypeContext) *ast.FuncType {
	ft := &ast.FuncType{
		IsAsync: ctx.ASYNC() != nil,
		Start:   t.pos(ctx.GetStart()),
	}
	if ctx.TypeList() != nil {
		for _, tr := range ctx.TypeList().AllTypeRef() {
			ft.Params = append(ft.Params, t.typeRef(tr))
		}
	}
	if ctx.TypeRef() != nil {
		ft.Results = []ast.TypeRef{t.typeRef(ctx.TypeRef())}
	}
	return ft
}

func (t *translator) baseType(ctx parser.IBaseTypeContext) ast.TypeRef {
	start := t.pos(ctx.GetStart())

	if ctx.PrimitiveType() != nil {
		return &ast.NamedType{Name: ctx.PrimitiveType().GetText(), TypePos: start}
	}
	// void, bool, string, byte, char — single keyword types
	for _, kw := range []antlr.TerminalNode{ctx.VOID(), ctx.BOOL(), ctx.STRING(), ctx.BYTE(), ctx.CHAR()} {
		if kw != nil {
			return &ast.NamedType{Name: kw.GetText(), TypePos: start}
		}
	}
	// vector[T]
	if ctx.VECTOR() != nil {
		typeRefs := ctx.AllTypeRef()
		return &ast.VectorType{Elem: t.typeRef(typeRefs[0]), Start: start}
	}
	// map[K]V  — two TypeRef children
	if ctx.MAP() != nil {
		typeRefs := ctx.AllTypeRef()
		return &ast.MapType{Key: t.typeRef(typeRefs[0]), Value: t.typeRef(typeRefs[1]), Start: start}
	}
	// []T (slice) — LBRACKET RBRACKET TypeRef, no Expression
	if ctx.LBRACKET() != nil && ctx.Expression() == nil {
		typeRefs := ctx.AllTypeRef()
		if len(typeRefs) == 1 {
			return &ast.SliceType{Elem: t.typeRef(typeRefs[0]), Start: start}
		}
	}
	// [N]T (fixed array) — LBRACKET Expression RBRACKET TypeRef
	if ctx.LBRACKET() != nil && ctx.Expression() != nil {
		typeRefs := ctx.AllTypeRef()
		return &ast.ArrayType{
			Len:   t.expr(ctx.Expression()),
			Elem:  t.typeRef(typeRefs[0]),
			Start: start,
		}
	}
	// Named / qualified, optionally generic
	if ctx.QualifiedName() != nil {
		var gargs []ast.TypeRef
		if ctx.GenericArgs() != nil {
			for _, tr := range ctx.GenericArgs().AllTypeRef() {
				gargs = append(gargs, t.typeRef(tr))
			}
		}
		return &ast.NamedType{Name: ctx.QualifiedName().GetText(), GenericArgs: gargs, TypePos: start}
	}
	return nil
}

// ─── Utilities ───────────────────────────────────────────────────────────────

func identTexts(nodes []antlr.TerminalNode) []string {
	out := make([]string, len(nodes))
	for i, n := range nodes {
		out[i] = n.GetText()
	}
	return out
}