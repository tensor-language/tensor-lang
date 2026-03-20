// ast/stmt.go
package ast

// BlockStmt is a braced list of statements.
type BlockStmt struct {
	List   []Stmt
	LBrace Position
	RBrace Position
}

func (s *BlockStmt) Pos() Position { return s.LBrace }
func (s *BlockStmt) End() Position { return s.RBrace }
func (s *BlockStmt) stmtNode()     {}

// DeclStmt wraps a VarDecl or ConstDecl when it appears inside a function body.
type DeclStmt struct {
	Decl Decl
}

func (s *DeclStmt) Pos() Position { return s.Decl.Pos() }
func (s *DeclStmt) End() Position { return s.Decl.End() }
func (s *DeclStmt) stmtNode()     {}

// AssignStmt covers all assignment forms: "x = y", "x += y", "x++".
// Op is "=", "+=", "-=", etc.; for "x++" it is "++" with no Rhs.
type AssignStmt struct {
	Target Expr
	Op     string
	Value  Expr // nil for increment/decrement
	Start  Position
}

func (s *AssignStmt) Pos() Position { return s.Start }
func (s *AssignStmt) End() Position {
	if s.Value != nil {
		return s.Value.End()
	}
	return s.Target.End()
}
func (s *AssignStmt) stmtNode() {}

// ReturnStmt represents "return expr" or "return (a, b)".
type ReturnStmt struct {
	Results []Expr // empty for bare "return"
	Start   Position
}

func (s *ReturnStmt) Pos() Position { return s.Start }
func (s *ReturnStmt) End() Position {
	if len(s.Results) > 0 {
		return s.Results[len(s.Results)-1].End()
	}
	return s.Start
}
func (s *ReturnStmt) stmtNode() {}

// BreakStmt represents "break".
type BreakStmt struct {
	Start Position
}

func (s *BreakStmt) Pos() Position { return s.Start }
func (s *BreakStmt) End() Position { return s.Start }
func (s *BreakStmt) stmtNode()     {}

// ContinueStmt represents "continue".
type ContinueStmt struct {
	Start Position
}

func (s *ContinueStmt) Pos() Position { return s.Start }
func (s *ContinueStmt) End() Position { return s.Start }
func (s *ContinueStmt) stmtNode()     {}

// DeferStmt represents "defer expr" — the expression is typically a CallExpr.
type DeferStmt struct {
	Call  Expr
	Start Position
}

func (s *DeferStmt) Pos() Position { return s.Start }
func (s *DeferStmt) End() Position { return s.Call.End() }
func (s *DeferStmt) stmtNode()     {}

// IfStmt represents the full if / else-if / else chain.
// Else is either another *IfStmt (else-if) or a *BlockStmt (else), or nil.
type IfStmt struct {
	Cond  Expr
	Body  *BlockStmt
	Else  Stmt
	Start Position
}

func (s *IfStmt) Pos() Position { return s.Start }
func (s *IfStmt) End() Position {
	if s.Else != nil {
		return s.Else.End()
	}
	return s.Body.End()
}
func (s *IfStmt) stmtNode() {}

// ForStmt covers C-style and while-style loops.
//   - C-style:     Init != nil, Cond != nil, Post != nil
//   - While-style: Init == nil, Cond != nil, Post == nil
//   - Infinite:    Init == nil, Cond == nil, Post == nil
type ForStmt struct {
	Init  Stmt // *DeclStmt or *AssignStmt; nil for while/infinite
	Cond  Expr // nil for infinite loop
	Post  Stmt // *AssignStmt; nil for while/infinite
	Body  *BlockStmt
	Start Position
}

func (s *ForStmt) Pos() Position { return s.Start }
func (s *ForStmt) End() Position { return s.Body.End() }
func (s *ForStmt) stmtNode()     {}

// ForInStmt covers "for item in collection" and "for k, v in map".
type ForInStmt struct {
	Key   string // always set
	Value string // empty for single-variable form
	Iter  Expr   // the collection or range expression
	Body  *BlockStmt
	Start Position
}

func (s *ForInStmt) Pos() Position { return s.Start }
func (s *ForInStmt) End() Position { return s.Body.End() }
func (s *ForInStmt) stmtNode()     {}

// SwitchCase is one "case X, Y: ..." branch.
type SwitchCase struct {
	Values []Expr // the comma-separated match values
	Body   []Stmt
	Start  Position
}

// SwitchStmt represents the full switch statement.
type SwitchStmt struct {
	Tag     Expr
	Cases   []*SwitchCase
	Default []Stmt // nil if no default branch
	Start   Position
	End_    Position
}

func (s *SwitchStmt) Pos() Position { return s.Start }
func (s *SwitchStmt) End() Position { return s.End_ }
func (s *SwitchStmt) stmtNode()     {}

// ExprStmt wraps a standalone expression used as a statement.
type ExprStmt struct {
	X Expr
}

func (s *ExprStmt) Pos() Position { return s.X.Pos() }
func (s *ExprStmt) End() Position { return s.X.End() }
func (s *ExprStmt) stmtNode()     {}