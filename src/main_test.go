package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/nickng/bibtex"
)

func mustParse(t *testing.T, s string) *bibtex.BibTex {
	t.Helper()
	bib, err := bibtex.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatalf("failed to parse bibtex: %v", err)
	}
	return bib
}

func TestRun(t *testing.T) {
	buf := bytes.NewBufferString("")
	bib := mustParse(t, `@inproceedings{Almutairi2024a,
		author = {Sultan Almutairi and Yogev Neumann and Khaled Harfoush},
		title = {Fingerprinting {VPNs} with Custom Router Firmware: A New Censorship Threat Model},
		booktitle = {Consumer Communications \& Networking Conference},
		publisher = {IEEE},
		year = {2024},
		url = {https://censorbib.nymity.ch/pdf/Almutairi2024a.pdf},
	}`)

	run(buf, bib)

	fmt.Println(buf)
}
