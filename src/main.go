package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/nickng/bibtex"
)

// Matches e.g.: @inproceedings{Doe2024a,
var re = regexp.MustCompile(`@[a-z]*\{([A-Za-z\-]*[0-9]{4}[a-z]),`)

// Map a cite name (e.g., Doe2024a) to its line number in the .bib file. All
// cite names are unique.
type entryToLineFunc func(string) int

// Augment bibtex.BibEntry with the entry's line number in the .bib file.
type bibEntry struct {
	bibtex.BibEntry
	lineNum int
}

func toStr(b bibtex.BibString) string {
	if b == nil {
		return ""
	}
	return b.String()
}

func parseBibFile(path string) []bibEntry {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	bib, err := bibtex.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	// Augment our BibTeX entries with their respective line numbers in the .bib
	// file. This is necessary to create the "Download BibTeX" links.
	lineOf := buildEntryToLineFunc(path)
	bibEntries := []bibEntry{}
	for _, entry := range bib.Entries {
		bibEntries = append(bibEntries, bibEntry{
			BibEntry: *entry,
			lineNum:  lineOf(entry.CiteName),
		})
	}

	return bibEntries
}

func buildEntryToLineFunc(path string) entryToLineFunc {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(file)
	entryToLine := make(map[string]int)
	line := 0
	for sc.Scan() {
		line++
		s := sc.Text()
		if !strings.HasPrefix(s, "@") {
			continue
		}
		entry := parseCiteName(s) // E.g., Doe2024a
		entryToLine[entry] = line
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
	}

	return func(entry string) int {
		if line, ok := entryToLine[entry]; ok {
			return line
		}
		log.Fatalf("could not find line number for cite name: %s", entry)
		return -1
	}
}

func parseCiteName(line string) string {
	matches := re.FindStringSubmatch(line)
	if len(matches) != 2 {
		log.Fatalf("failed to extract cite name of: %s", line)
	}
	return matches[1]
}

func run(w io.Writer, bibEntries []bibEntry) {
	fmt.Fprint(w, header())
	makeBib(w, bibEntries)
	fmt.Fprint(w, footer())
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
