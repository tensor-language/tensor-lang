// ast/decl.go
package ast

// File is the root node for a parsed source file.
type File struct {
	Namespace string
	Imports   []*ImportDecl
	Decls     []Decl
	Start     Position
}

func (f *File) Pos() Position { return f.Start }
func (f *File) End() Position {
	if len(f.Decls) == 0 {
		return f.Start
	}
	return f.Decls[len(f.Decls)-1].End()
}

// Attribute represents @packed or @align(16).
type Attribute struct {
	Name  string
	Arg   Expr // nil if no argument
	Start Position
}

func (a *Attribute) Pos() Position { return a.Start }
func (a *Attribute) End() Position { return a.Start }

// ImportDecl represents a single import path with an optional alias.
// Alias is "_" for blank imports, "." for dot imports, or a name; empty means no alias.
type ImportDecl struct {
	Alias string
	Path  string
	Start Position
}

func (d *ImportDecl) Pos() Position { return d.Start }
func (d *ImportDecl) End() Position { return d.Start }
func (d *ImportDecl) declNode()     {}

// ConstSpec is a single name = value pair inside a const declaration.
type ConstSpec struct {
	Name  string
	Type  TypeRef // may be nil (inferred)
	Value Expr
	Start Position
}

func (s *ConstSpec) Pos() Position { return s.Start }
func (s *ConstSpec) End() Position { return s.Start }

// ConstDecl represents "const X = 1" or "const ( X = 1; Y = 2 )".
type ConstDecl struct {
	Specs []*ConstSpec
	Start Position
}

func (d *ConstDecl) Pos() Position { return d.Start }
func (d *ConstDecl) End() Position { return d.Start }
func (d *ConstDecl) declNode()     {}

// VarDecl represents "var x: T = val" or "let x = val".
// IsRef == true means var (heap ref-counted); false means let (stack).
type VarDecl struct {
	Name   string
	IsRef  bool
	Type   TypeRef // may be nil (inferred)
	Value  Expr    // nil for "var x: T = null" when null is the value
	IsNull bool    // true when the initialiser is the literal null
	Start  Position
}

func (d *VarDecl) Pos() Position { return d.Start }
func (d *VarDecl) End() Position { return d.Start }
func (d *VarDecl) declNode()     {}

// SelfParam describes the receiver on a method or deinit.
type SelfParam struct {
	Name  string
	Type  TypeRef
	IsMut bool // true for "&mut self"
	Start Position
}

func (s *SelfParam) Pos() Position { return s.Start }
func (s *SelfParam) End() Position { return s.Start }

// Field is a name:type pair used in params, interface fields, and generic params.
type Field struct {
	Name  string
	Type  TypeRef
	Start Position
}

func (f *Field) Pos() Position { return f.Start }
func (f *Field) End() Position { return f.Start }

// FuncDecl represents any function: plain, async, gpu, method, or generic.
type FuncDecl struct {
	Attrs         []*Attribute
	Name          string
	IsAsync       bool
	IsGpu         bool
	Self          *SelfParam // nil for free functions
	GenericParams []string   // ["T", "U"] — names only; constraints TBD
	Params        []*Field
	IsVariadic    bool       // true if last param is "..."
	ReturnType    TypeRef    // nil for void; TupleType for multiple returns
	Body          *BlockStmt // nil for extern declarations
	Start         Position
}

func (d *FuncDecl) Pos() Position { return d.Start }
func (d *FuncDecl) End() Position {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Start
}
func (d *FuncDecl) declNode() {}

// DeinitDecl represents "deinit(self x: T) { ... }".
type DeinitDecl struct {
	Self  *SelfParam
	Body  *BlockStmt
	Start Position
}

func (d *DeinitDecl) Pos() Position { return d.Start }
func (d *DeinitDecl) End() Position { return d.Body.End() }
func (d *DeinitDecl) declNode()     {}

// InterfaceDecl represents "interface Foo[T] { ... }".
type InterfaceDecl struct {
	Attrs         []*Attribute
	Name          string
	GenericParams []string
	Fields        []*Field
	Start         Position
}

func (d *InterfaceDecl) Pos() Position { return d.Start }
func (d *InterfaceDecl) End() Position { return d.Start }
func (d *InterfaceDecl) declNode()     {}

// EnumMember is a single member of an enum.
type EnumMember struct {
	Name  string
	Value Expr // nil if auto-assigned
	Start Position
}

func (m *EnumMember) Pos() Position { return m.Start }
func (m *EnumMember) End() Position { return m.Start }

// EnumDecl represents "enum Direction: uint8 { ... }".
type EnumDecl struct {
	Name           string
	UnderlyingType TypeRef // nil → int32 default
	Members        []*EnumMember
	Start          Position
}

func (d *EnumDecl) Pos() Position { return d.Start }
func (d *EnumDecl) End() Position { return d.Start }
func (d *EnumDecl) declNode()     {}

// TypeAliasDecl represents "type FILE = opaque" or "type X = SomeType".
type TypeAliasDecl struct {
	Name     string
	IsOpaque bool
	Type     TypeRef // nil when IsOpaque
	Start    Position
}

func (d *TypeAliasDecl) Pos() Position { return d.Start }
func (d *TypeAliasDecl) End() Position { return d.Start }
func (d *TypeAliasDecl) declNode()     {}

// ExternDecl represents "extern c { ... }" or "extern cpp { ... }".
// The members are kept as raw data for the frontend to lower into
// symbol-table entries; we don't need a full sub-AST for them yet.
type ExternDecl struct {
	Lang    string // "c" or "cpp"
	Members []ExternMember
	Start   Position
}

func (d *ExternDecl) Pos() Position { return d.Start }
func (d *ExternDecl) End() Position { return d.Start }
func (d *ExternDecl) declNode()     {}

// ExternMember is a sealed interface over the kinds of things
// that can appear inside an extern block.
type ExternMember interface {
	externMemberNode()
}

// ExternFunc describes a single C/C++ function binding.
type ExternFunc struct {
	Convention string  // "cdecl", "stdcall", etc.; empty = default
	Name       string  // Arc-side name
	Symbol     string  // C-side name if different; empty = same as Name
	Params     []TypeRef
	IsVariadic bool
	Return     TypeRef // nil = void
	Start      Position
}

func (e *ExternFunc) Pos() Position    { return e.Start }
func (e *ExternFunc) End() Position    { return e.Start }
func (e *ExternFunc) externMemberNode() {}

// ExternTypeAlias is "type Comparator = func(*void,*void) int32" inside extern.
type ExternTypeAlias struct {
	Name  string
	Type  TypeRef
	Start Position
}

func (e *ExternTypeAlias) Pos() Position    { return e.Start }
func (e *ExternTypeAlias) End() Position    { return e.Start }
func (e *ExternTypeAlias) externMemberNode() {}

// ExternType describes a C/C++ struct, union, or typedef binding inside an extern block.
// Kind is one of: "struct", "union", "typedef".
type ExternType struct {
	Kind   string   // "struct", "union", or "typedef"
	Name   string   // Arc-side name
	Symbol string   // C-side name if different; empty = same as Name
	Fields []*Field // nil for opaque / typedef forms
	Start  Position
}

func (e *ExternType) Pos() Position    { return e.Start }
func (e *ExternType) End() Position    { return e.Start }
func (e *ExternType) externMemberNode() {}

// ExternNamespace wraps a "namespace Foo { ... }" block inside extern cpp.
type ExternNamespace struct {
	Name    string
	Members []ExternMember
	Start   Position
}

func (e *ExternNamespace) Pos() Position    { return e.Start }
func (e *ExternNamespace) End() Position    { return e.Start }
func (e *ExternNamespace) externMemberNode() {}

// ExternClass wraps a "class ID3D11Device { ... }" block inside extern cpp.
type ExternClass struct {
	IsAbstract bool
	Name       string
	Symbol     string
	Methods    []*ExternMethod
	Start      Position
}

func (e *ExternClass) Pos() Position    { return e.Start }
func (e *ExternClass) End() Position    { return e.Start }
func (e *ExternClass) externMemberNode() {}

// ExternMethodKind distinguishes the four kinds of class members.
type ExternMethodKind int

const (
	ExternVirtual ExternMethodKind = iota
	ExternStatic
	ExternConstructor
	ExternDestructor
)

// ExternMethod is any method inside an extern class.
type ExternMethod struct {
	Kind       ExternMethodKind
	Convention string
	Name       string
	Symbol     string
	Params     []TypeRef // includes self as first param for virtual/destructor
	IsVariadic bool
	Return     TypeRef
	IsConst    bool // true for "const" methods
	Start      Position
}

func (e *ExternMethod) Pos() Position { return e.Start }
func (e *ExternMethod) End() Position { return e.Start }