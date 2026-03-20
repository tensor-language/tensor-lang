// ast/ast.go
package ast

type Position struct {
	Line   int
	Column int
}

type Node interface {
	Pos() Position
	End() Position
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Decl interface {
	Node
	declNode()
}

// Visitor is implemented by any pass that walks the AST.
// Return false from any Visit* to skip that node's children.
type Visitor interface {
	VisitFile(*File) bool
	VisitFuncDecl(*FuncDecl) bool
	VisitDeinitDecl(*DeinitDecl) bool
	VisitInterfaceDecl(*InterfaceDecl) bool
	VisitEnumDecl(*EnumDecl) bool
	VisitConstDecl(*ConstDecl) bool
	VisitVarDecl(*VarDecl) bool
	VisitTypeAliasDecl(*TypeAliasDecl) bool
	VisitExternDecl(*ExternDecl) bool
	VisitImportDecl(*ImportDecl) bool

	VisitBlockStmt(*BlockStmt) bool
	VisitIfStmt(*IfStmt) bool
	VisitForStmt(*ForStmt) bool
	VisitForInStmt(*ForInStmt) bool
	VisitSwitchStmt(*SwitchStmt) bool
	VisitReturnStmt(*ReturnStmt) bool
	VisitDeferStmt(*DeferStmt) bool
	VisitAssignStmt(*AssignStmt) bool
	VisitBreakStmt(*BreakStmt) bool
	VisitContinueStmt(*ContinueStmt) bool
	VisitExprStmt(*ExprStmt) bool
	VisitDeclStmt(*DeclStmt) bool

	VisitBinaryExpr(*BinaryExpr) bool
	VisitUnaryExpr(*UnaryExpr) bool
	VisitCallExpr(*CallExpr) bool
	VisitSelectorExpr(*SelectorExpr) bool
	VisitIndexExpr(*IndexExpr) bool
	VisitSliceExpr(*SliceExpr) bool
	VisitRangeExpr(*RangeExpr) bool
	VisitAwaitExpr(*AwaitExpr) bool
	VisitIdent(*Ident) bool
	VisitBasicLit(*BasicLit) bool
	VisitCompositeLit(*CompositeLit) bool
	VisitKeyValueExpr(*KeyValueExpr) bool
	VisitLambdaExpr(*LambdaExpr) bool
	VisitTupleLit(*TupleLit) bool
	VisitProcessExpr(*ProcessExpr) bool
	VisitNewExpr(*NewExpr) bool
	VisitNewArrayExpr(*NewArrayExpr) bool
	VisitDeleteExpr(*DeleteExpr) bool
	VisitCastExpr(*CastExpr) bool
}