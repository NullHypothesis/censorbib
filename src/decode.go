package main

import (
	"log"
	"strings"
)

type conversion struct {
	from string
	to   string
}

func decodeTitle(title string) string {
	for _, convert := range []conversion{
		{`\#`, "#"},
		{`--`, `–`},
		{"``", "“"},
		{"''", "”"},
		{"'", "’"},       // U+2019
		{`$\cdot$`, `·`}, // U+00B7.
	} {
		title = strings.ReplaceAll(title, convert.from, convert.to)
	}

	// Get rid of all curly brackets. We're displaying titles without changing
	// their casing.
	title = strings.ReplaceAll(title, "{", "")
	title = strings.ReplaceAll(title, "}", "")

	return title
}

func decodeAuthors(authors string) string {
	for _, convert := range []conversion{
		{"'", "’"},
	} {
		authors = strings.ReplaceAll(authors, convert.from, convert.to)
	}
	// For simplicity, we expect authors to be formatted as "John Doe" instead
	// of "Doe, John".
	if strings.Contains(authors, ",") {
		log.Fatalf("author %q contains a comma", authors)
	}
	authorSlice := strings.Split(authors, " and ")
	return strings.Join(authorSlice, ", ")
}

func decodeProceedings(proceedings string) string {
	for _, convert := range []conversion{
		{`\&`, "&"},
	} {
		proceedings = strings.ReplaceAll(proceedings, convert.from, convert.to)
	}
	return proceedings
}
