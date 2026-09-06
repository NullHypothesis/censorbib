package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/nickng/bibtex"
)

func TestUniqueBibTeXKeys(t *testing.T) {
	file, err := os.Open("../references.bib")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	bib, err := bibtex.Parse(file)
	if err != nil {
		t.Fatalf("failed to parse references.bib: %v", err)
	}

	seen := make(map[string]int)
	for i, entry := range bib.Entries {
		if first, ok := seen[entry.CiteName]; ok {
			t.Errorf("duplicate BibTeX key %q in references.bib (entries %d and %d)", entry.CiteName, first, i+1)
		} else {
			seen[entry.CiteName] = i + 1
		}
	}
}

func TestBibTeXIndentation(t *testing.T) {
	contents, err := os.ReadFile("../references.bib")
	if err != nil {
		t.Fatal(err)
	}

	for i, line := range bytes.Split(contents, []byte("\n")) {
		for _, char := range line {
			if char == ' ' {
				t.Errorf("references.bib:%d: use tabs instead of spaces for indentation", i+1)
				break
			}
			if char != '\t' {
				break
			}
		}
	}
}
