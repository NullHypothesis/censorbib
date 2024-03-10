package main

import (
	"fmt"
	"io"
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

func appendIfNotEmpty(slice []string, s string) []string {
	if s != "" {
		return append(slice, s)
	}
	return slice
}

func makeBib(to io.Writer, bibEntries []bibEntry) {
	yearToEntries := make(map[string][]string)

	for _, entry := range bibEntries {
		y := entry.Fields["year"].String()
		yearToEntries[y] = append(yearToEntries[y], makeBibEntry(&entry))
	}

	sortedYears := sortByYear(yearToEntries)
	for _, year := range sortedYears {
		fmt.Fprintln(to, "<ul>")
		for _, entry := range yearToEntries[year] {
			fmt.Fprint(to, entry)
		}
		fmt.Fprintln(to, "</ul>")
	}
}

func makeBibEntry(entry *bibEntry) string {
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

func makeBibEntryTitle(entry *bibEntry) string {
	// Paper title is on the left side.
	title := []string{
		`<span class="paper">`,
		decodeTitle(entry.Fields["title"].String()),
		`</span>`,
	}
	// Icons are on the right side.
	icons := []string{
		`<span class="icons">`,
		fmt.Sprintf("<a href='%s'>", entry.Fields["url"].String()),
		`<img class="icon" title="Download paper" src="assets/pdf-icon.svg" alt="Download icon">`,
		`</a>`,
		fmt.Sprintf("<a href='https://censorbib.nymity.ch/pdf/%s.pdf'>", entry.CiteName),
		`<img class="icon" title="Download cached paper" src="assets/cache-icon.svg" alt="Cached download icon">`,
		`</a>`,
		fmt.Sprintf("<a href='https://github.com/NullHypothesis/censorbib/blob/master/references.bib#L%d'>", entry.lineNum),
		`<img class="icon" title="Download BibTeX" src="assets/bibtex-icon.svg" alt="BibTeX download icon">`,
		`</a>`,
		fmt.Sprintf("<a href='#%s'>", entry.CiteName),
		`<img class="icon" title="Link to paper" src="assets/link-icon.svg" alt="Paper link icon">`,
		`</a>`,
		`</span>`,
	}
	return strings.Join(append(title, icons...), "\n")
}

func makeBibEntryAuthors(entry *bibEntry) string {
	s := []string{
		`<span class="author">`,
		decodeAuthors(entry.Fields["author"].String()),
		`</span>`,
	}
	return strings.Join(s, "\n")
}

func makeBibEntryMisc(entry *bibEntry) string {
	s := []string{}
	s = appendIfNotEmpty(s, makeBibEntryVenue(entry))
	s = appendIfNotEmpty(s, toStr(entry.Fields["year"]))
	s = appendIfNotEmpty(s, toStr(entry.Fields["publisher"]))
	return strings.Join(s, ", ")
}

func makeBibEntryVenue(entry *bibEntry) string {
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
		return "" // Some entries are self-published.
	}

	s := []string{
		prefix,
		`<span class="venue">`,
		decodeProceedings(toStr(bs)),
		`</span>`,
	}

	return strings.Join(s, "")
}
