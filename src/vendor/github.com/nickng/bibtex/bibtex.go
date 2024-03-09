package bibtex

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

// BibString is a segment of a bib string.
type BibString interface {
	RawString() string // Internal representation.
	String() string    // Displayed string.
}

// BibVar is a string variable.
type BibVar struct {
	Key   string    // Variable key.
	Value BibString // Variable actual value.
}

// RawString is the internal representation of the variable.
func (v *BibVar) RawString() string {
	return v.Key
}

func (v *BibVar) String() string {
	return v.Value.String()
}

// BibConst is a string constant.
type BibConst string

// NewBibConst converts a constant string to BibConst.
func NewBibConst(c string) BibConst {
	return BibConst(c)
}

// RawString is the internal representation of the constant (i.e. the string).
func (c BibConst) RawString() string {
	return fmt.Sprintf("{%s}", string(c))
}

func (c BibConst) String() string {
	return string(c)
}

// BibComposite is a composite string, may contain both variable and string.
type BibComposite []BibString

// NewBibComposite creates a new composite with one element.
func NewBibComposite(s BibString) *BibComposite {
	comp := &BibComposite{}
	return comp.Append(s)
}

// Append adds a BibString to the composite
func (c *BibComposite) Append(s BibString) *BibComposite {
	comp := append(*c, s)
	return &comp
}

func (c *BibComposite) String() string {
	var buf bytes.Buffer
	for _, s := range *c {
		buf.WriteString(s.String())
	}
	return buf.String()
}

// RawString returns a raw (bibtex) representation of the composite string.
func (c *BibComposite) RawString() string {
	var buf bytes.Buffer
	for i, comp := range *c {
		if i > 0 {
			buf.WriteString(" # ")
		}
		switch comp := comp.(type) {
		case *BibConst:
			buf.WriteString(comp.RawString())
		case *BibVar:
			buf.WriteString(comp.RawString())
		case *BibComposite:
			buf.WriteString(comp.RawString())
		}
	}
	return buf.String()
}

// BibEntry is a record of BibTeX record.
type BibEntry struct {
	Type     string
	CiteName string
	Fields   map[string]BibString
}

// NewBibEntry creates a new BibTeX entry.
func NewBibEntry(entryType string, citeName string) *BibEntry {
	spaceStripper := strings.NewReplacer(" ", "")
	cleanedType := strings.ToLower(spaceStripper.Replace(entryType))
	cleanedName := spaceStripper.Replace(citeName)
	return &BibEntry{
		Type:     cleanedType,
		CiteName: cleanedName,
		Fields:   map[string]BibString{},
	}
}

// AddField adds a field (key-value) to a BibTeX entry.
func (entry *BibEntry) AddField(name string, value BibString) {
	entry.Fields[strings.TrimSpace(name)] = value
}

// prettyStringConfig controls the formatting/printing behaviour of the BibTex's and BibEntry's PrettyPrint functions
type prettyStringConfig struct {
	// priority controls the order in which fields are printed. Keys with lower values are printed earlier.
	//See keyOrderToPriorityMap
	priority map[string]int
}

// keyOrderToPriorityMap is a helper function for WithKeyOrder, converting the user facing key order slice
// into the map format that is internally used by the sort function
func keyOrderToPriorityMap(keyOrder []string) map[string]int {
	priority := make(map[string]int)
	offset := len(keyOrder)
	for i, v := range keyOrder {
		priority[v] = i - offset
	}
	return priority
}

var defaultPrettyStringConfig = prettyStringConfig{priority: keyOrderToPriorityMap([]string{"title", "author", "url"})}

// PrettyStringOpt allows to change the pretty print format for BibEntry and BibTex
type PrettyStringOpt func(config *prettyStringConfig)

// WithKeyOrder changes the order in which BibEntry keys are printed to the order in which they appear in keyOrder
func WithKeyOrder(keyOrder []string) PrettyStringOpt {
	return func(config *prettyStringConfig) {
		config.priority = make(map[string]int)
		offset := len(keyOrder)
		for i, v := range keyOrder {
			config.priority[v] = i - offset
		}
	}
}

// prettyStringAppend appends the pretty print string for BibEntry using config to configure the formatting
func (entry *BibEntry) prettyStringAppend(buf *bytes.Buffer, config prettyStringConfig) {
	fmt.Fprintf(buf, "@%s{%s,\n", entry.Type, entry.CiteName)

	// Determine key order.
	keys := []string{}
	for key := range entry.Fields {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		pi, pj := config.priority[keys[i]], config.priority[keys[j]]
		return pi < pj || (pi == pj && keys[i] < keys[j])
	})

	// Write fields.
	tw := tabwriter.NewWriter(buf, 1, 4, 1, ' ', 0)
	for _, key := range keys {
		value := entry.Fields[key].String()
		format := stringformat(value)
		fmt.Fprintf(tw, "    %s\t=\t"+format+",\n", key, value)
	}
	tw.Flush()
	buf.WriteString("}\n")

}

// PrettyString pretty prints a BibEntry
func (entry *BibEntry) PrettyString(options ...PrettyStringOpt) string {
	config := defaultPrettyStringConfig
	for _, option := range options {
		option(&config)
	}
	var buf bytes.Buffer
	entry.prettyStringAppend(&buf, config)

	return buf.String()
}

// String returns a BibTex entry as a simplified BibTex string.
func (entry *BibEntry) String() string {
	var bibtex bytes.Buffer
	bibtex.WriteString(fmt.Sprintf("@%s{%s,\n", entry.Type, entry.CiteName))
	for key, val := range entry.Fields {
		if i, err := strconv.Atoi(strings.TrimSpace(val.String())); err == nil {
			bibtex.WriteString(fmt.Sprintf("  %s = %d,\n", key, i))
		} else {
			bibtex.WriteString(fmt.Sprintf("  %s = {%s},\n", key, strings.TrimSpace(val.String())))
		}
	}
	bibtex.Truncate(bibtex.Len() - 2)
	bibtex.WriteString(fmt.Sprintf("\n}\n"))
	return bibtex.String()
}

// RawString returns a BibTex entry data structure in its internal representation.
func (entry *BibEntry) RawString() string {
	var bibtex bytes.Buffer
	bibtex.WriteString(fmt.Sprintf("@%s{%s,\n", entry.Type, entry.CiteName))
	for key, val := range entry.Fields {
		if i, err := strconv.Atoi(strings.TrimSpace(val.String())); err == nil {
			bibtex.WriteString(fmt.Sprintf("  %s = %d,\n", key, i))
		} else {
			bibtex.WriteString(fmt.Sprintf("  %s = %s,\n", key, val.RawString()))
		}
	}
	bibtex.Truncate(bibtex.Len() - 2)
	bibtex.WriteString(fmt.Sprintf("\n}\n"))
	return bibtex.String()
}

// BibTex is a list of BibTeX entries.
type BibTex struct {
	Preambles []BibString        // List of Preambles
	Entries   []*BibEntry        // Items in a bibliography.
	StringVar map[string]*BibVar // Map from string variable to string.

	// A list of default BibVars that are implicitly
	// defined and can be used without defining
	defaultVars map[string]string
}

// NewBibTex creates a new BibTex data structure.
func NewBibTex() *BibTex {
	// Sets up some default vars
	months := map[string]time.Month{
		"jan": 1, "feb": 2, "mar": 3,
		"apr": 4, "may": 5, "jun": 6,
		"jul": 7, "aug": 8, "sep": 9,
		"oct": 10, "nov": 11, "dec": 12,
	}

	defaultVars := make(map[string]string)
	for mth, month := range months {
		// TODO(nickng): i10n of month name in user's local language
		defaultVars[mth] = month.String()
	}

	return &BibTex{
		Preambles: []BibString{},
		Entries:   []*BibEntry{},
		StringVar: make(map[string]*BibVar),

		defaultVars: defaultVars,
	}
}

// AddPreamble adds a preamble to a bibtex.
func (bib *BibTex) AddPreamble(p BibString) {
	bib.Preambles = append(bib.Preambles, p)
}

// AddEntry adds an entry to the BibTeX data structure.
func (bib *BibTex) AddEntry(entry *BibEntry) {
	bib.Entries = append(bib.Entries, entry)
}

// AddStringVar adds a new string var (if does not exist).
func (bib *BibTex) AddStringVar(key string, val BibString) {
	bib.StringVar[key] = &BibVar{Key: key, Value: val}
}

// GetStringVar looks up a string by its key.
func (bib *BibTex) GetStringVar(key string) *BibVar {
	if bv, ok := bib.StringVar[key]; ok {
		return bv
	}
	if v, ok := bib.getDefaultVar(key); ok {
		return v
	}
	// This is undefined.
	log.Fatalf("%s: %s", ErrUnknownStringVar, key)
	return nil
}

// getDefaultVar is a fallback for looking up keys (e.g. 3-character month)
// and use them even though it hasn't been defined in the bib.
func (bib *BibTex) getDefaultVar(key string) (*BibVar, bool) {
	if v, ok := bib.defaultVars[key]; ok {
		// if found, add this to the BibTex
		bib.StringVar[key] = &BibVar{Key: key, Value: NewBibConst(v)}
		return bib.StringVar[key], true
	}

	return nil, false
}

// String returns a BibTex data structure as a simplified BibTex string.
func (bib *BibTex) String() string {
	var bibtex bytes.Buffer
	for _, entry := range bib.Entries {
		bibtex.WriteString(entry.String())
	}
	return bibtex.String()
}

// RawString returns a BibTex data structure in its internal representation.
func (bib *BibTex) RawString() string {
	var bibtex bytes.Buffer
	for k, strvar := range bib.StringVar {
		bibtex.WriteString(fmt.Sprintf("@string{%s = {%s}}\n", k, strvar.String()))
	}
	for _, preamble := range bib.Preambles {
		bibtex.WriteString(fmt.Sprintf("@preamble{%s}\n", preamble.RawString()))
	}
	for _, entry := range bib.Entries {
		bibtex.WriteString(entry.RawString())
	}
	return bibtex.String()
}

// PrettyString pretty prints a BibTex
func (bib *BibTex) PrettyString(options ...PrettyStringOpt) string {
	config := defaultPrettyStringConfig
	for _, option := range options {
		option(&config)
	}

	var buf bytes.Buffer
	for i, entry := range bib.Entries {
		if i != 0 {
			fmt.Fprint(&buf, "\n")
		}
		entry.prettyStringAppend(&buf, config)

	}
	return buf.String()
}

// stringformat determines the correct formatting verb for the given BibTeX field value.
func stringformat(v string) string {
	// Numbers may be represented unquoted.
	if _, err := strconv.Atoi(v); err == nil {
		return "%s"
	}

	// Strings with certain characters must be brace quoted.
	if strings.ContainsAny(v, "\"{}") {
		return "{%s}"
	}

	// Default to quoted string.
	return "%q"
}
