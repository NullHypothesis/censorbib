package bibtex

import (
	"fmt"
	"strings"
)

// Lexer token.
type token int

const (
	// tILLEGAL stands for an invalid token.
	tILLEGAL token = iota
)

var eof = rune(0)

// tokenPos is a pair of coordinate to identify start of token.
type tokenPos struct {
	Char  int
	Lines []int
}

func (p tokenPos) String() string {
	return fmt.Sprintf("%d:%d", len(p.Lines)+1, p.Char)
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isAlpha(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9')
}

func isAlphanum(ch rune) bool {
	return isAlpha(ch) || isDigit(ch)
}

func isBareSymbol(ch rune) bool {
	return strings.ContainsRune("-_:./+", ch)
}

// isSymbol returns true if ch is a valid symbol
func isSymbol(ch rune) bool {
	return strings.ContainsRune("!?&*+-./:;<>[]^_`|~@", ch)
}

func isOpenQuote(ch rune) bool {
	return ch == '{' || ch == '"'
}
