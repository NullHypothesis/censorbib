package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nickng/bibtex"
)

func mustParse(t *testing.T, s string) bibEntry {
	t.Helper()
	bib, err := bibtex.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatalf("failed to parse bibtex: %v", err)
	}
	return bibEntry{
		BibEntry: *bib.Entries[0],
		lineNum:  0,
	}
}

func TestRun(t *testing.T) {
	buf := bytes.NewBufferString("")
	entry := mustParse(t, `@inproceedings{Almutairi2024a,
		author = {Sultan Almutairi and Yogev Neumann and Khaled Harfoush},
		title = {Fingerprinting {VPNs} with Custom Router Firmware: A New Censorship Threat Model},
		booktitle = {Consumer Communications \& Networking Conference},
		publisher = {IEEE},
		year = {2024},
		url = {https://censorbib.nymity.ch/pdf/Almutairi2024a.pdf},
	}`)

	makeBib(buf, []bibEntry{entry})

	bufStr := buf.String()
	if !strings.HasPrefix(bufStr, "<ul>") {
		t.Errorf("expected <ul> but got %q...", bufStr[:10])
	}
	if !strings.HasSuffix(bufStr, "</ul>\n") {
		t.Errorf("expected </ul> but got %q", bufStr[len(bufStr)-10:])
	}
}
