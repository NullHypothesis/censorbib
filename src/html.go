package main

import (
	"bytes"
	"html/template"
	"io"
	"sort"
	"strings"

	"github.com/nickng/bibtex"
)

type bibEntryView struct {
	CiteName      string
	Title         string
	Authors       string
	Venue         string
	VenuePrefix   string
	HasVenue      bool
	Year          string
	HasYear       bool
	Publisher     string
	URL           string
	DiscussionURL string
}

var bibEntryTemplate = template.Must(template.New("bib-entry").Parse(`<li id="{{.CiteName}}">
<div>
<span class="paper">{{.Title}}</span>
<span class="icons">
{{if .DiscussionURL}}<a href="{{.DiscussionURL}}"><img class="icon" title="Online discussion" src="assets/discussion-icon.svg" alt="Discussion icon"></a>{{end}}
<a href="{{.URL}}"><img class="icon" title="Download paper" src="assets/pdf-icon.svg" alt="Download icon"></a>
<a href="https://censorbib-papers.t3.tigrisfiles.io/{{.CiteName}}.pdf"><img class="icon" title="Download cached paper" src="assets/cache-icon.svg" alt="Cached download icon"></a>
<a href="#bibtex-{{.CiteName}}" class="bibtex-link" data-reference="{{.CiteName}}" title="Show BibTeX" aria-label="Show BibTeX for {{.Title}}"><img class="icon" src="assets/bibtex-icon.svg" alt="BibTeX icon"></a>
<a href="#{{.CiteName}}"><img class="icon" title="Link to paper" src="assets/link-icon.svg" alt="Paper link icon"></a>
</span>
</div>
<div>
<span class="author">{{.Authors}}</span>
</div>
<span class="other">{{if .HasVenue}}{{.VenuePrefix}}<span class="venue">{{.Venue}}</span>{{end}}{{if .Year}}{{if .HasVenue}}, {{end}}{{.Year}}{{end}}{{if .Publisher}}{{if or .HasVenue .HasYear}}, {{end}}{{.Publisher}}{{end}}</span>
</li>
`))

func makeBib(to io.Writer, bibEntries []bibEntry) {
	previousYear := ""
	for _, entry := range bibEntries {
		year := toStr(entry.Fields["year"])
		if year != previousYear {
			if previousYear != "" {
				mustFprintln(to, "</ul>")
			}
			mustFprintf(to, "<ul class=\"year-group\" data-year=\"%s\">\n", template.HTMLEscapeString(year))
			previousYear = year
		}
		mustFprint(to, makeBibEntry(&entry))
	}
	if previousYear != "" {
		mustFprintln(to, "</ul>")
	}
}

func makeBibEntry(entry *bibEntry) string {
	buf := new(bytes.Buffer)
	if err := bibEntryTemplate.Execute(buf, entryView(entry)); err != nil {
		panic(err)
	}
	return buf.String()
}

func entryView(entry *bibEntry) bibEntryView {
	prefix, venue := entryVenueParts(entry)
	year := toStr(entry.Fields["year"])
	return bibEntryView{
		CiteName:      entry.CiteName,
		Title:         entryTitle(entry),
		Authors:       entryAuthors(entry),
		Venue:         venue,
		VenuePrefix:   prefix,
		HasVenue:      venue != "",
		Year:          year,
		HasYear:       year != "",
		Publisher:     toStr(entry.Fields["publisher"]),
		URL:           toStr(entry.Fields["url"]),
		DiscussionURL: toStr(entry.Fields["discussion_url"]),
	}
}

func entryTitle(entry *bibEntry) string {
	return decodeTitle(toStr(entry.Fields["title"]))
}

func entryAuthors(entry *bibEntry) string {
	return decodeAuthors(toStr(entry.Fields["author"]))
}

func entryVenue(entry *bibEntry) string {
	_, venue := entryVenueParts(entry)
	return venue
}

func entryVenueParts(entry *bibEntry) (string, string) {
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
		return "", "" // Some entries are self-published.
	}

	return prefix, decodeProceedings(toStr(bs))
}

func sortBibEntries(bibEntries []bibEntry) {
	sort.SliceStable(bibEntries, func(i, j int) bool {
		a := &bibEntries[i]
		b := &bibEntries[j]
		for _, cmp := range []struct {
			left       string
			right      string
			descending bool
		}{
			{toStr(a.Fields["year"]), toStr(b.Fields["year"]), true},
			{entryVenue(a), entryVenue(b), false},
			{entryTitle(a), entryTitle(b), false},
			{a.CiteName, b.CiteName, false},
		} {
			left := strings.ToLower(cmp.left)
			right := strings.ToLower(cmp.right)
			if left == right {
				continue
			}
			if cmp.descending {
				return left > right
			}
			return left < right
		}
		return false
	})
}

func makeSearchBox(to io.Writer, count int) {
	mustFprintf(to, `<form id="search-form" role="search" action="">
  <label for="search-input">Search</label>
  <input id="search-input" type="search" name="q" autocomplete="off" placeholder="Title, author, venue, year, publisher, or cite name">
  <span id="result-count" aria-live="polite">%d papers</span>
</form>
<div id="no-results" hidden>No matches.</div>
`, count)
}
