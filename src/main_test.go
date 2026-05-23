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
		BibEntry:  *bib.Entries[0],
		rawBibtex: strings.TrimSpace(s),
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
	if !strings.HasPrefix(bufStr, `<ul class="year-group"`) {
		t.Errorf("expected year group <ul> but got %q...", bufStr[:30])
	}
	if !strings.HasSuffix(bufStr, "</ul>\n") {
		t.Errorf("expected </ul> but got %q", bufStr[len(bufStr)-10:])
	}
}

func TestExtractRawBibEntries(t *testing.T) {
	contents := `@inproceedings{Doe2024a,
	author = {Jane Doe},
	title = {A paper with {nested} braces},
	year = {2024},
}

@article{Müller2023a,
	author = {Max Müller},
	title = {Another paper},
	year = {2023}
}`

	rawByCiteName := extractRawBibEntries([]byte(contents))
	if got := rawByCiteName["Doe2024a"]; !strings.Contains(got, `{nested}`) {
		t.Fatalf("raw BibTeX did not preserve nested braces: %q", got)
	}
	if got := rawByCiteName["Müller2023a"]; !strings.HasPrefix(got, "@article{Müller2023a") {
		t.Fatalf("raw BibTeX did not preserve non-ASCII cite name: %q", got)
	}
}

func TestSortBibEntries(t *testing.T) {
	entries := []bibEntry{
		mustParse(t, `@inproceedings{Beta2024a,
			author = {Jane Doe},
			title = {Beta},
			booktitle = {Workshop},
			year = {2024},
			url = {https://example.com/beta.pdf},
		}`),
		mustParse(t, `@inproceedings{Alpha2025a,
			author = {Jane Doe},
			title = {Alpha},
			booktitle = {Workshop},
			year = {2025},
			url = {https://example.com/alpha.pdf},
		}`),
		mustParse(t, `@inproceedings{Zeta2024a,
			author = {Jane Doe},
			title = {Zeta},
			booktitle = {Conference},
			year = {2024},
			url = {https://example.com/zeta.pdf},
		}`),
	}

	sortBibEntries(entries)
	got := []string{entries[0].CiteName, entries[1].CiteName, entries[2].CiteName}
	want := []string{"Alpha2025a", "Zeta2024a", "Beta2024a"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("unexpected sort order: got %v, want %v", got, want)
	}
}

func TestMakeReferenceDataScript(t *testing.T) {
	buf := bytes.NewBufferString("")
	entry := mustParse(t, `@inproceedings{Doe2024a,
		author = {Jane Doe},
		title = {Searchable Paper},
		booktitle = {Free and Open Communications on the Internet},
		publisher = {Example Publisher},
		year = {2024},
		url = {https://example.com/paper.pdf},
	}`)

	makeReferenceDataScript(buf, []bibEntry{entry})
	got := buf.String()
	for _, want := range []string{`"citeName":"Doe2024a"`, `"title":"Searchable Paper"`, `"publisher":"Example Publisher"`, `"rawBibtex":"@inproceedings{Doe2024a,`} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated metadata missing %q in %s", want, got)
		}
	}
}
