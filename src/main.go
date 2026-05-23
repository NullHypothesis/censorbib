package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/nickng/bibtex"
)

// Matches e.g.: @inproceedings{Müller2024a,
var re = regexp.MustCompile(`(?i)^@[a-z]+\s*\{\s*([^,\s]+)\s*,`)

// Augment bibtex.BibEntry with the entry's raw record in the .bib file.
type bibEntry struct {
	bibtex.BibEntry
	rawBibtex string
}

type searchEntry struct {
	CiteName  string `json:"citeName"`
	Title     string `json:"title"`
	Authors   string `json:"authors"`
	Venue     string `json:"venue"`
	Year      string `json:"year"`
	Publisher string `json:"publisher"`
	RawBibtex string `json:"rawBibtex"`
}

func toStr(b bibtex.BibString) string {
	if b == nil {
		return ""
	}
	return b.String()
}

func parseBibFile(path string) []bibEntry {
	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	file := bytes.NewReader(contents)
	bib, err := bibtex.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	rawByCiteName := extractRawBibEntries(contents)
	bibEntries := []bibEntry{}
	for _, entry := range bib.Entries {
		rawBibtex, ok := rawByCiteName[entry.CiteName]
		if !ok {
			log.Fatalf("could not find raw BibTeX for cite name: %s", entry.CiteName)
		}
		bibEntries = append(bibEntries, bibEntry{
			BibEntry:  *entry,
			rawBibtex: rawBibtex,
		})
	}

	return bibEntries
}

func extractRawBibEntries(contents []byte) map[string]string {
	rawByCiteName := make(map[string]string)
	for i := 0; i < len(contents); i++ {
		if contents[i] != '@' {
			continue
		}

		depth := 0
		for j := i; j < len(contents); j++ {
			switch contents[j] {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					raw := strings.TrimSpace(string(contents[i : j+1]))
					rawByCiteName[parseCiteName(raw)] = raw
					i = j
					goto nextEntry
				}
			}
		}
	nextEntry:
	}
	return rawByCiteName
}

func parseCiteName(line string) string {
	matches := re.FindStringSubmatch(line)
	if len(matches) != 2 {
		log.Fatalf("failed to extract cite name of: %s", line)
	}
	return matches[1]
}

func mustFprint(w io.Writer, a ...any) {
	if _, err := fmt.Fprint(w, a...); err != nil {
		log.Fatalf("failed to write HTML: %v", err)
	}
}

func mustFprintln(w io.Writer, a ...any) {
	if _, err := fmt.Fprintln(w, a...); err != nil {
		log.Fatalf("failed to write HTML: %v", err)
	}
}

func mustFprintf(w io.Writer, format string, a ...any) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		log.Fatalf("failed to write HTML: %v", err)
	}
}

func run(w io.Writer, bibEntries []bibEntry) {
	sortBibEntries(bibEntries)
	mustFprint(w, header())
	makeSearchBox(w, len(bibEntries))
	mustFprintln(w, "<div id='container'>")
	makeBib(w, bibEntries)
	mustFprintln(w, "</div>")
	makeReferenceDataScript(w, bibEntries)
	mustFprint(w, footer())
}

func makeReferenceDataScript(w io.Writer, bibEntries []bibEntry) {
	searchEntries := []searchEntry{}
	for _, entry := range bibEntries {
		searchEntries = append(searchEntries, searchEntry{
			CiteName:  entry.CiteName,
			Title:     entryTitle(&entry),
			Authors:   entryAuthors(&entry),
			Venue:     entryVenue(&entry),
			Year:      toStr(entry.Fields["year"]),
			Publisher: toStr(entry.Fields["publisher"]),
			RawBibtex: entry.rawBibtex,
		})
	}

	mustFprintln(w, `<script id="reference-data" type="application/json">`)
	if err := json.NewEncoder(w).Encode(searchEntries); err != nil {
		log.Fatalf("failed to encode search data: %v", err)
	}
	mustFprintln(w, `</script>`)
}

func main() {
	path := flag.String("path", "", "Path to .bib file.")
	flag.Parse()
	if *path == "" {
		log.Fatal("No path to .bib file provided.")
	}
	run(os.Stdout, parseBibFile(*path))
	log.Println("Successfully created bibliography.")
}
