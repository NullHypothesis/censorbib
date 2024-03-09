package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/nickng/bibtex"
)

func sortByYear(yearToEntries map[string][]string) []string {
	keys := make([]string, 0, len(yearToEntries))
	for k := range yearToEntries {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	return keys
}

func makeBib(to io.Writer, bib *bibtex.BibTex) {
	yearToEntries := make(map[string][]string)

	for _, entry := range bib.Entries {
		y := entry.Fields["year"].String()
		yearToEntries[y] = append(yearToEntries[y], makeBibEntry(entry))
	}

	sortedYears := sortByYear(yearToEntries)
	for _, year := range sortedYears {
		fmt.Fprintf(to, "<ul>")
		for _, entry := range yearToEntries[year] {
			fmt.Fprint(to, entry)
		}
		fmt.Fprintf(to, "</ul>")
	}
}

func makeBibEntry(entry *bibtex.BibEntry) string {
	s := []string{
		fmt.Sprintf("<li id='%s'>", entry.CiteName),
		`<div>`,
		makeBibEntryTitle(entry),
		`</div>`,
		`<div>`,
		makeBibEntryAuthors(entry),
		`</div>`,
		`<span class="other">`,
		makeBibEntryMisc(entry),
		`</span>`,
		`</li>`,
	}
	return strings.Join(s, "\n")
}

func makeBibEntryTitle(entry *bibtex.BibEntry) string {

	// Paper title is on the left side.
	title := []string{
		`<span class="paper">`,
		//fmt.Sprintf("<a name='%s'>", entry.CiteName),
		decodeTitle(entry.Fields["title"].String()),
		//`</a>`,
		`</span>`,
	}
	// Icons are on the right side.
	icons := []string{
		`<span class="icons">`,
		fmt.Sprintf("<a href='%s'>", entry.Fields["url"].String()),
		`<img class="icon" title="Download paper" src="img/pdf-icon.svg" alt="Download icon">`,
		`</a>`,
		fmt.Sprintf("<a href='pdf/%s.pdf'>", entry.CiteName),
		`<img class="icon" title="Download cached paper" src="img/cache-icon.svg" alt="Cached download icon">`,
		`</a>`,
		fmt.Sprintf("<a href='bibtex.html#%s'>", entry.CiteName),
		`<img class="icon" title="Download BibTeX" src="img/bibtex-icon.svg" alt="BibTeX download icon">`,
		`</a>`,
		fmt.Sprintf("<a href='#%s'>", entry.CiteName),
		`<img class="icon" title="Link to paper" src="img/link-icon.svg" alt="Paper link icon">`,
		`</a>`,
		`</span>`,
	}
	return strings.Join(append(title, icons...), "\n")
}

func makeBibEntryAuthors(entry *bibtex.BibEntry) string {
	s := []string{
		`<span class="author">`,
		decodeAuthors(entry.Fields["author"].String()),
		`</span>`,
	}
	return strings.Join(s, "\n")
}

func makeBibEntryMisc(entry *bibtex.BibEntry) string {
	s := []string{}
	s = appendIfNotEmpty(s, makeBibEntryVenue(entry))
	s = appendIfNotEmpty(s, toStr(entry.Fields["year"]))
	s = appendIfNotEmpty(s, toStr(entry.Fields["publisher"]))
	return strings.Join(s, ", ")
}

func makeBibEntryVenue(entry *bibtex.BibEntry) string {
	var (
		prefix string
		bs     bibtex.BibString
		ok     bool
	)

	if bs, ok = entry.Fields["booktitle"]; ok {
		prefix = "In Proc. of: "
	} else if bs, ok = entry.Fields["journal"]; ok {
		prefix = "In: "
	} else {
		// Some entries are self-published.
		return ""
	}

	s := []string{
		prefix,
		`<span class="venue">`,
		decodeProceedings(toStr(bs)),
		`</span>`,
	}

	return strings.Join(s, "")
}

func appendIfNotEmpty(slice []string, s string) []string {
	if s != "" {
		return append(slice, s)
	}
	return slice
}

func toStr(b bibtex.BibString) string {
	if b == nil {
		return ""
	}
	return b.String()
}

func run(w io.Writer, b *bibtex.BibTex) {
	fmt.Fprint(w, header())
	makeBib(w, b)
	fmt.Fprint(w, footer())
}

func parseBibFile(path string) *bibtex.BibTex {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	bib, err := bibtex.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	return bib
}

func main() {
	path := flag.String("path", "", "Path to .bib file.")
	flag.Parse()
	if *path == "" {
		log.Fatal("No path to .bib file provided.")
	}
	run(os.Stdout, parseBibFile(*path))
}
