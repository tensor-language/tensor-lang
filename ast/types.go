// ast/types.go
package ast

// TypeRef is the sealed interface for all type-reference nodes.
type TypeRef interface {
	Node
	typeNode()
}

// NamedType represents "int32", "string", "Point", "Box[int32]", "net.Socket".
type NamedType struct {
	Name        string // may contain a single dot for qualified names: "net.Socket"
	GenericArgs []TypeRef
	TypePos     Position
}

func (t *NamedType) Pos() Position { return t.TypePos }
func (t *NamedType) End() Position { return t.TypePos }
func (t *NamedType) typeNode()     {}

// PointerType represents "*T" — only valid inside extern blocks.
type PointerType struct {
	IsConst bool // *const T
	Base    TypeRef
	Start   Position
}

func (t *PointerType) Pos() Position { return t.Start }
func (t *PointerType) End() Position { return t.Base.End() }
func (t *PointerType) typeNode()     {}

// RefType represents "&T" and "&const T" — only valid inside extern cpp blocks.
type RefType struct {
	IsConst bool
	Base    TypeRef
	Start   Position
}

func (t *RefType) Pos() Position { return t.Start }
func (t *RefType) End() Position { return t.Base.End() }
func (t *RefType) typeNode()     {}

// MutRefType represents "&mut T" in regular Arc parameter position.
type MutRefType struct {
	Base  TypeRef
	Start Position
}

func (t *MutRefType) Pos() Position { return t.Start }
func (t *MutRefType) End() Position { return t.Base.End() }
func (t *MutRefType) typeNode()     {}

// SliceType represents "[]T".
type SliceType struct {
	Elem  TypeRef
	Start Position
}

func (t *SliceType) Pos() Position { return t.Start }
func (t *SliceType) End() Position { return t.Elem.End() }
func (t *SliceType) typeNode()     {}

// ArrayType represents "[N]T" with a compile-time constant length.
type ArrayType struct {
	Len   Expr
	Elem  TypeRef
	Start Position
}

func (t *ArrayType) Pos() Position { return t.Start }
func (t *ArrayType) End() Position { return t.Elem.End() }
func (t *ArrayType) typeNode()     {}

// VectorType represents "vector[T]".
type VectorType struct {
	Elem  TypeRef
	Start Position
}

func (t *VectorType) Pos() Position { return t.Start }
func (t *VectorType) End() Position { return t.Elem.End() }
func (t *VectorType) typeNode()     {}

// MapType represents "map[K]V".
type MapType struct {
	Key   TypeRef
	Value TypeRef
	Start Position
}

func (t *MapType) Pos() Position { return t.Start }
func (t *MapType) End() Position { return t.Value.End() }
func (t *MapType) typeNode()     {}

// FuncType represents "func(T, U) V" or "async func(T) V" as a type (not a declaration).
type FuncType struct {
	IsAsync bool
	Params  []TypeRef
	Results []TypeRef // empty = void; length > 1 = tuple return
	Start   Position
}

func (t *FuncType) Pos() Position { return t.Start }
func (t *FuncType) End() Position { return t.Start }
func (t *FuncType) typeNode()     {}

// TupleType represents "(T, U)" in return-type position only.
type TupleType struct {
	Elems []TypeRef
	Start Position
}

func (t *TupleType) Pos() Position { return t.Start }
func (t *TupleType) End() Position { return t.Start }
func (t *TupleType) typeNode()     {}