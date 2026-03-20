package syntax

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
)

// SyntaxError represents a specific parsing error with location information.
type SyntaxError struct {
	Line    int
	Column  int
	Message string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("syntax error at line %d:%d: %s", e.Line, e.Column, e.Message)
}

// ErrorListener implements the ANTLR ErrorListener interface to collect errors.
type ErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []*SyntaxError
}

func NewErrorListener() *ErrorListener {
	return &ErrorListener{
		Errors: make([]*SyntaxError, 0),
	}
}

// SyntaxError is the callback fired by ANTLR when parsing fails.
func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.Errors = append(l.Errors, &SyntaxError{
		Line:    line,
		Column:  column,
		Message: msg,
	})
}