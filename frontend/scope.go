// frontend/scope.go
package frontend

import "github.com/arc-language/arc-lang/ast"

// ScopeKind classifies what created a scope.
type ScopeKind int

const (
	ScopeGlobal ScopeKind = iota
	ScopeFile
	ScopeFunc
	ScopeBlock
)

// Symbol is a named entity resolved during semantic analysis.
type Symbol struct {
	Name string
	Kind string   // "var", "let", "const", "func", "type", "param", "enum", "enumMember"
	Decl ast.Node // The AST node that introduced this name
	Type ast.TypeRef // Filled in during type-checking
}

// Scope is a lexical scope that chains to its parent.
type Scope struct {
	Parent  *Scope
	Kind    ScopeKind
	Symbols map[string]*Symbol
}

func NewScope(parent *Scope, kind ScopeKind) *Scope {
	return &Scope{
		Parent:  parent,
		Kind:    kind,
		Symbols: make(map[string]*Symbol),
	}
}

// Insert adds a symbol, returning any previous symbol with the same name
// so the caller can emit a "redeclared" error if needed.
func (s *Scope) Insert(name string, sym *Symbol) *Symbol {
	prev := s.Symbols[name]
	s.Symbols[name] = sym
	return prev
}

// Lookup walks the scope chain, returning nil if the name is not found.
func (s *Scope) Lookup(name string) *Symbol {
	if sym, ok := s.Symbols[name]; ok {
		return sym
	}
	if s.Parent != nil {
		return s.Parent.Lookup(name)
	}
	return nil
}

// LookupLocal checks only this scope, not parents.
func (s *Scope) LookupLocal(name string) *Symbol {
	return s.Symbols[name]
}

// Enclosing returns the nearest enclosing scope of the given kind, or nil.
func (s *Scope) Enclosing(kind ScopeKind) *Scope {
	cur := s
	for cur != nil {
		if cur.Kind == kind {
			return cur
		}
		cur = cur.Parent
	}
	return nil
}