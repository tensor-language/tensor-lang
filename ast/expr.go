// ast/expr.go
package ast

// BinaryExpr represents "x + y".
type BinaryExpr struct {
	Left   Expr
	Op     string // "+", "-", "*", "/", "&&", "||", "&", "|", etc.
	Right  Expr
	OpPos  Position
}

func (e *BinaryExpr) Pos() Position { return e.Left.Pos() }
func (e *BinaryExpr) End() Position { return e.Right.End() }
func (e *BinaryExpr) exprNode()     {}

// UnaryExpr represents "-x", "!x", "~x".
type UnaryExpr struct {
	Op    string // "-", "!", "~"
	X     Expr
	Start Position
}

func (e *UnaryExpr) Pos() Position { return e.Start }
func (e *UnaryExpr) End() Position { return e.X.End() }
func (e *UnaryExpr) exprNode()     {}

// SelectorExpr represents "x.field".
type SelectorExpr struct {
	X    Expr
	Sel  string
	Dot  Position
}

func (e *SelectorExpr) Pos() Position { return e.X.Pos() }
func (e *SelectorExpr) End() Position { return e.Dot }
func (e *SelectorExpr) exprNode()     {}

// CallExpr represents "fn(args)" and "obj.method(args)".
type CallExpr struct {
	Fun    Expr
	Args   []Expr
	Start  Position
	EndPos Position
}

func (e *CallExpr) Pos() Position { return e.Start }
func (e *CallExpr) End() Position { return e.EndPos }
func (e *CallExpr) exprNode()     {}

// AwaitExpr represents "await expr".
type AwaitExpr struct {
	X     Expr
	Start Position
}

func (e *AwaitExpr) Pos() Position { return e.Start }
func (e *AwaitExpr) End() Position { return e.X.End() }
func (e *AwaitExpr) exprNode()     {}

// Ident represents a bare name "x".
type Ident struct {
	Name    string
	NamePos Position
}

func (e *Ident) Pos() Position { return e.NamePos }
func (e *Ident) End() Position { return e.NamePos }
func (e *Ident) exprNode()     {}

// BasicLit represents a literal "123", "3.14", `"hello"`, "'a'".
// Kind is one of: INT, HEX, FLOAT, STRING, CHAR, BOOL, NULL.
type BasicLit struct {
	Kind    string
	Value   string
	LitPos  Position
}

func (e *BasicLit) Pos() Position { return e.LitPos }
func (e *BasicLit) End() Position { return e.LitPos }
func (e *BasicLit) exprNode()     {}

// CompositeLit represents "Point{x: 1, y: 2}" or "vector[int32]{1,2,3}".
// Type may be nil for bare "{}" initializers inferred from context.
type CompositeLit struct {
	Type   TypeRef // nil for anonymous initializers
	Fields []Expr  // KeyValueExpr entries or plain Expr for slices/vectors
	LBrace Position
	RBrace Position
}

func (e *CompositeLit) Pos() Position { return e.LBrace }
func (e *CompositeLit) End() Position { return e.RBrace }
func (e *CompositeLit) exprNode()     {}

// KeyValueExpr represents "key: value" inside composite literals.
type KeyValueExpr struct {
	Key   Expr
	Value Expr
	Colon Position
}

func (e *KeyValueExpr) Pos() Position { return e.Key.Pos() }
func (e *KeyValueExpr) End() Position { return e.Value.End() }
func (e *KeyValueExpr) exprNode()     {}

// IndexExpr represents "arr[i]".
type IndexExpr struct {
	X      Expr
	Index  Expr
	LBrack Position
	RBrack Position
}

func (e *IndexExpr) Pos() Position { return e.X.Pos() }
func (e *IndexExpr) End() Position { return e.RBrack }
func (e *IndexExpr) exprNode()     {}

// SliceExpr represents "buf[lo..hi]".
type SliceExpr struct {
	X      Expr
	Low    Expr
	High   Expr
	LBrack Position
	RBrack Position
}

func (e *SliceExpr) Pos() Position { return e.X.Pos() }
func (e *SliceExpr) End() Position { return e.RBrack }
func (e *SliceExpr) exprNode()     {}

// RangeExpr represents "0..10" (appears in for-in headers and slice bounds).
type RangeExpr struct {
	Low   Expr
	High  Expr
	DotDot Position
}

func (e *RangeExpr) Pos() Position { return e.Low.Pos() }
func (e *RangeExpr) End() Position { return e.High.End() }
func (e *RangeExpr) exprNode()     {}

// LambdaExpr represents "(x: int32) => { ... }" and "async (x: int32) => { ... }".
type LambdaExpr struct {
	IsAsync bool
	Params  []*Field
	Body    *BlockStmt
	Start   Position
}

func (e *LambdaExpr) Pos() Position { return e.Start }
func (e *LambdaExpr) End() Position { return e.Body.End() }
func (e *LambdaExpr) exprNode()     {}

// TupleLit represents "(a, b)" in return position.
type TupleLit struct {
	Elems  []Expr
	LParen Position
	RParen Position
}

func (e *TupleLit) Pos() Position { return e.LParen }
func (e *TupleLit) End() Position { return e.RParen }
func (e *TupleLit) exprNode()     {}

// ProcessExpr represents "process func(...) { ... }(args)".
type ProcessExpr struct {
	Params []*Field
	Body   *BlockStmt
	Args   []Expr
	Start  Position
}

func (e *ProcessExpr) Pos() Position { return e.Start }
func (e *ProcessExpr) End() Position { return e.Body.End() }
func (e *ProcessExpr) exprNode()     {}

// NewExpr represents "new Point{x:1}".
type NewExpr struct {
	Type TypeRef
	Init *CompositeLit
	Start Position
}

func (e *NewExpr) Pos() Position { return e.Start }
func (e *NewExpr) End() Position { return e.Init.End() }
func (e *NewExpr) exprNode()     {}

// NewArrayExpr represents "new [4096]byte".
type NewArrayExpr struct {
	Len   Expr
	Elem  TypeRef
	Start Position
}

func (e *NewArrayExpr) Pos() Position { return e.Start }
func (e *NewArrayExpr) End() Position { return e.Elem.End() }
func (e *NewArrayExpr) exprNode()     {}

// DeleteExpr represents "delete(ptr)" — used in defer and bare statements.
type DeleteExpr struct {
	X     Expr
	Start Position
	End_  Position
}

func (e *DeleteExpr) Pos() Position { return e.Start }
func (e *DeleteExpr) End() Position { return e.End_ }
func (e *DeleteExpr) exprNode()     {}

// CastExpr represents "int32(x)" — type-name in call position.
// The frontend distinguishes a cast from a regular call by checking
// whether Fun resolves to a type name rather than a value.
type CastExpr struct {
	Type  TypeRef
	X     Expr
	Start Position
	End_  Position
}

func (e *CastExpr) Pos() Position { return e.Start }
func (e *CastExpr) End() Position { return e.End_ }
func (e *CastExpr) exprNode()     {}