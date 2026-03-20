package syntax

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/arc-language/arc-lang/parser"
)

// tokenSliceSource implements antlr.TokenSource by embedding the original lexer.
// Embedding *parser.ArcLexer ensures we satisfy ALL interface requirements
// (including private methods like setTokenFactory) automatically.
type tokenSliceSource struct {
	*parser.ArcLexer
	tokens []antlr.Token
	index  int
}

// NextToken overrides the embedded lexer's method to serve tokens from our slice.
func (s *tokenSliceSource) NextToken() antlr.Token {
	if s.index >= len(s.tokens) {
		return antlr.CommonTokenFactoryDEFAULT.Create(
			s.GetTokenSourceCharStreamPair(),
			antlr.TokenEOF,
			"<EOF>",
			antlr.TokenDefaultChannel,
			-1, -1, 0, 0,
		)
	}
	t := s.tokens[s.index]
	s.index++
	return t
}

// createTokenStream wraps the generated lexer to perform Automatic Semicolon Insertion (ASI).
// It returns a token stream that includes synthetic SEMI tokens where newlines imply termination.
func createTokenStream(input string) *antlr.CommonTokenStream {
	inputStream := antlr.NewInputStream(input)
	lexer := parser.NewArcLexer(inputStream)
	lexer.RemoveErrorListeners()

	// 1. Fetch all raw tokens from the ANTLR lexer (includes hidden channel tokens).
	allTokens := lexer.GetAllTokens()
	var processedTokens []antlr.Token

	// Grab the lexer's source pair once for minting synthetic tokens.
	sourcePair := lexer.GetTokenSourceCharStreamPair()

	// 2. Iterate through tokens and inject SEMI where appropriate.
	for i, t := range allTokens {
		processedTokens = append(processedTokens, t)

		if shouldInsertSemi(t, i, allTokens) {
			semi := antlr.CommonTokenFactoryDEFAULT.Create(
				sourcePair,
				parser.ArcLexerSEMI,
				";",
				antlr.TokenDefaultChannel,
				-1, -1,
				t.GetLine(),
				t.GetColumn()+len(t.GetText()),
			)
			processedTokens = append(processedTokens, semi)
		}
	}

	// 3. Create a wrapper that embeds the original lexer but serves the processed tokens.
	wrapper := &tokenSliceSource{
		ArcLexer: lexer,
		tokens:   processedTokens,
		index:    0,
	}

	// 4. Create the stream using our wrapper as the source.
	return antlr.NewCommonTokenStream(wrapper, antlr.TokenDefaultChannel)
}

// shouldInsertSemi determines if a SEMI token should be inserted after token t.
// It skips hidden-channel tokens (NL, comments, whitespace) when looking for
// the next visible token, so that line comparisons are made against real tokens only.
func shouldInsertSemi(t antlr.Token, index int, all []antlr.Token) bool {
	// Only consider tokens on the default channel as trigger candidates.
	if t.GetChannel() != antlr.TokenDefaultChannel {
		return false
	}

	// Find the next token on the default channel, skipping hidden ones.
	var next antlr.Token
	for j := index + 1; j < len(all); j++ {
		if all[j].GetChannel() == antlr.TokenDefaultChannel {
			next = all[j]
			break
		}
	}

	// If the next visible token is on the same line, no semi needed.
	if next != nil && next.GetLine() == t.GetLine() {
		return false
	}

	// next == nil means EOF; otherwise next is on a later line.
	// Either way, check whether this token type triggers insertion.
	switch t.GetTokenType() {
	case parser.ArcLexerIDENTIFIER,
		parser.ArcLexerINT_LIT,
		parser.ArcLexerHEX_LIT,
		parser.ArcLexerFLOAT_LIT,
		parser.ArcLexerSTRING_LIT,
		parser.ArcLexerCHAR_LIT,
		parser.ArcLexerTRUE,
		parser.ArcLexerFALSE,
		parser.ArcLexerNULL,
		parser.ArcLexerRETURN,
		parser.ArcLexerBREAK,
		parser.ArcLexerCONTINUE,
		parser.ArcLexerRPAREN,
		parser.ArcLexerRBRACKET,
		parser.ArcLexerRBRACE,
		parser.ArcLexerINC,
		parser.ArcLexerDEC:
		return true
	case parser.ArcLexerINT8, parser.ArcLexerINT16, parser.ArcLexerINT32, parser.ArcLexerINT64,
		parser.ArcLexerUINT8, parser.ArcLexerUINT16, parser.ArcLexerUINT32, parser.ArcLexerUINT64,
		parser.ArcLexerFLOAT32, parser.ArcLexerFLOAT64,
		parser.ArcLexerBOOL, parser.ArcLexerSTRING, parser.ArcLexerBYTE, parser.ArcLexerCHAR,
		parser.ArcLexerUSIZE, parser.ArcLexerISIZE:
		return true
	}
	return false
}