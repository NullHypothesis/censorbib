package bibtex

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

var parseField bool

// scanner is a lexical scanner
type scanner struct {
	commentMode bool
	r           *bufio.Reader
	pos         tokenPos
}

// newScanner returns a new instance of scanner.
func newScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r), pos: tokenPos{Char: 0, Lines: []int{}}}
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.eof is returned).
func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	if ch == '\n' {
		s.pos.Lines = append(s.pos.Lines, s.pos.Char)
		s.pos.Char = 0
	} else {
		s.pos.Char++
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *scanner) unread() {
	_ = s.r.UnreadRune()
	if s.pos.Char == 0 {
		s.pos.Char = s.pos.Lines[len(s.pos.Lines)-1]
		s.pos.Lines = s.pos.Lines[:len(s.pos.Lines)-1]
	} else {
		s.pos.Char--
	}
}

// Scan returns the next token and literal value.
func (s *scanner) Scan() (tok token, lit string, err error) {
	ch := s.read()
	if isWhitespace(ch) {
		s.ignoreWhitespace()
		ch = s.read()
	}
	if isAlphanum(ch) {
		s.unread()
		return s.scanIdent()
	}
	switch ch {
	case eof:
		return 0, "", nil
	case '@':
		return tATSIGN, string(ch), nil
	case ':':
		return tCOLON, string(ch), nil
	case ',':
		parseField = false // reset parseField if reached end of field.
		return tCOMMA, string(ch), nil
	case '=':
		parseField = true // set parseField if = sign outside quoted or ident.
		return tEQUAL, string(ch), nil
	case '"':
		tok, lit := s.scanQuoted()
		return tok, lit, nil
	case '{':
		if parseField {
			return s.scanBraced()
		}
		// If we're reading a comment, return everything after {
		// to the next @-sign (exclusive)
		if s.commentMode {
			s.unread()
			commentBodyTok, commentBody := s.scanCommentBody()
			return commentBodyTok, commentBody, nil
		}
		return tLBRACE, string(ch), nil
	case '}':
		if parseField { // reset parseField if reached end of entry.
			parseField = false
		}
		return tRBRACE, string(ch), nil
	case '#':
		return tPOUND, string(ch), nil
	case ' ':
		s.ignoreWhitespace()
	}
	return tILLEGAL, string(ch), nil
}

// scanIdent categorises a string to one of three categories.
func (s *scanner) scanIdent() (tok token, lit string, err error) {
	switch ch := s.read(); ch {
	case '"':
		tok, lit := s.scanQuoted()
		return tok, lit, nil
	case '{':
		return s.scanBraced()
	default:
		s.unread() // Not open quote/brace.
		tok, lit := s.scanBare()
		return tok, lit, nil
	}
}

func (s *scanner) scanBare() (token, string) {
	var buf bytes.Buffer
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isAlphanum(ch) && !isBareSymbol(ch) || isWhitespace(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	str := buf.String()
	if strings.ToLower(str) == "comment" {
		s.commentMode = true
		return tCOMMENT, str
	} else if strings.ToLower(str) == "preamble" {
		return tPREAMBLE, str
	} else if strings.ToLower(str) == "string" {
		return tSTRING, str
	} else if _, err := strconv.Atoi(str); err == nil && parseField { // Special case for numeric
		return tIDENT, str
	}
	return tBAREIDENT, str
}

// scanBraced parses a braced string, like {this}.
func (s *scanner) scanBraced() (token, string, error) {
	var buf bytes.Buffer
	var macro bool
	brace := 1
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '\\' {
			_, _ = buf.WriteRune(ch)
			macro = true
		} else if ch == '{' {
			_, _ = buf.WriteRune(ch)
			brace++
		} else if ch == '}' {
			brace--
			macro = false
			if brace == 0 { // Balances open brace.
				return tIDENT, buf.String(), nil
			}
			_, _ = buf.WriteRune(ch)
		} else if ch == '@' {
			if macro {
				_, _ = buf.WriteRune(ch)
			} else {
				return token(0), buf.String(), ErrUnexpectedAtsign
			}
		} else if isWhitespace(ch) {
			_, _ = buf.WriteRune(ch)
			macro = false
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	return tILLEGAL, buf.String(), nil
}

// scanQuoted parses a quoted string, like "this".
func (s *scanner) scanQuoted() (token, string) {
	var buf bytes.Buffer
	brace := 0
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '{' {
			brace++
		} else if ch == '}' {
			brace--
		} else if ch == '"' {
			if brace == 0 { // Matches open quote, unescaped
				return tIDENT, buf.String()
			}
			_, _ = buf.WriteRune(ch)
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	return tILLEGAL, buf.String()
}

// skipCommentBody is a scan method used for reading bibtex
// comment item by reading all runes until the next @.
//
// e.g.
// @comment{...anything can go here even if braces are unbalanced@
// comment body string will be "...anything can go here even if braces are unbalanced"
func (s *scanner) scanCommentBody() (token, string) {
	var buf bytes.Buffer
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '@' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	s.commentMode = false
	return tCOMMENTBODY, buf.String()
}

// ignoreWhitespace consumes the current rune and all contiguous whitespace.
func (s *scanner) ignoreWhitespace() {
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
	}
}
