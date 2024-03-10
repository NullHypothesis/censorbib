//go:generate goyacc -p bibtex -o bibtex.y.go bibtex.y

package bibtex

import (
	"fmt"
	"io"
)

// lexer for bibtex.
type lexer struct {
	scanner     *scanner
	ParseErrors []error // Parse errors from yacc
	Errors      []error // Other errors
}

// newLexer returns a new yacc-compatible lexer.
func newLexer(r io.Reader) *lexer {
	return &lexer{
		scanner: newScanner(r),
	}
}

// Lex is provided for yacc-compatible parser.
func (l *lexer) Lex(yylval *bibtexSymType) int {
	token, strval, err := l.scanner.Scan()
	if err != nil {
		l.Errors = append(l.Errors, fmt.Errorf("%w at %s", err, l.scanner.pos))
		return int(0)
	}
	yylval.strval = strval
	return int(token)
}

// Error handles error.
func (l *lexer) Error(err string) {
	l.ParseErrors = append(l.ParseErrors, &ErrParse{Err: err, Pos: l.scanner.pos})
}
