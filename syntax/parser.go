package syntax

import (
	"github.com/arc-language/arc-lang/parser"
)

// ParseResult contains the root of the Concrete Syntax Tree (CST) and any syntax errors.
type ParseResult struct {
	Root   parser.ICompilationUnitContext
	Errors []*SyntaxError
}

// Parse takes a string of Arc source code and returns the parsed CST.
// It automatically handles lexing and semicolon insertion.
func Parse(input string) *ParseResult {
	// 1. Lexing with Automatic Semicolon Insertion (ASI)
	stream := createTokenStream(input)

	// 2. Initialize the ANTLR Parser
	p := parser.NewArcParser(stream)

	// 3. Attach Custom Error Listener
	// We remove the default console listener to prevent noise.
	p.RemoveErrorListeners()
	listener := NewErrorListener()
	p.AddErrorListener(listener)

	// 4. Execute the Parse
	// CompilationUnit is the top-level rule defined in ArcParser.g4
	tree := p.CompilationUnit()

	return &ParseResult{
		Root:   tree,
		Errors: listener.Errors,
	}
}