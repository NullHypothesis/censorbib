// Package bibtex is a bibtex parser written in Go.
//
// The package contains a simple parser and data structure to represent bibtex
// records.
//
// # Supported syntax
//
// The basic syntax is:
//
//	@BIBTYPE{IDENT,
//	    key1 = word,
//	    key2 = "quoted",
//	    key3 = {quoted},
//	}
//
// where BIBTYPE is the type of document (e.g. inproceedings, article, etc.)
// and IDENT is a string identifier.
//
// The bibtex format is not standardised, this parser follows the descriptions
// found in the link below. If there are any problems, please file any issues
// with a minimal working example at the GitHub repository.
// http://maverick.inria.fr/~Xavier.Decoret/resources/xdkbibtex/bibtex_summary.html
package bibtex // import "github.com/nickng/bibtex"
