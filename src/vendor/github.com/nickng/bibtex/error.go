package bibtex

import (
	"errors"
	"fmt"
)

var (
	// ErrUnexpectedAtsign is an error for unexpected @ in {}.
	ErrUnexpectedAtsign = errors.New("unexpected @ sign")
	// ErrUnknownStringVar is an error for looking up undefined string var.
	ErrUnknownStringVar = errors.New("unknown string variable")
)

// ErrParse is a parse error.
type ErrParse struct {
	Pos tokenPos
	Err string // Error string returned from parser.
}

func (e *ErrParse) Error() string {
	return fmt.Sprintf("parse failed at %s: %s", e.Pos, e.Err)
}
