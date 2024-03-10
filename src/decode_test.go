package main

import (
	"testing"
)

func TestToString(t *testing.T) {
	testCases := []conversion{
		{
			from: "Title",
			to:   "Title",
		},
		{
			from: "This is a {Title}",
			to:   "This is a Title",
		},
		{
			from: "This is a {Title}",
			to:   "This is a Title",
		},
		{
			from: `{\#h00t}: Censorship Resistant Microblogging`,
			to:   `#h00t: Censorship Resistant Microblogging`,
		},
		{
			from: "``Good'' Worms and Human Rights",
			to:   "“Good” Worms and Human Rights",
		},
		{
			from: "An Analysis of {China}'s ``{Great Cannon}''",
			to:   "An Analysis of China’s “Great Cannon”",
		},
		{
			from: `lib$\cdot$erate, (n):`,
			to:   `lib·erate, (n):`,
		},
		{
			from: "Well -- Exploring the {Great} {Firewall}'s Poisoned {DNS}",
			to:   "Well – Exploring the Great Firewall’s Poisoned DNS",
		},
	}

	for _, test := range testCases {
		to := decodeTitle(test.from)
		if to != test.to {
			t.Errorf("Expected\n%s\ngot\n%s", test.to, to)
		}
	}
}

func TestDecodeAuthors(t *testing.T) {
	testCases := []conversion{
		{ // Multiple authors should be separated by commas.
			from: "John Doe and Jane Doe",
			to:   "John Doe, Jane Doe",
		},
		{ // Single authors should remain as-is.
			from: "John Doe",
			to:   "John Doe",
		},
		{ // Single-name authors should remain as-is.
			from: "John and Jane",
			to:   "John, Jane",
		},
		{ // Non-ASCII characters should be unaffected.
			from: "Jóhn Doe",
			to:   "Jóhn Doe",
		},
		{ // Apostrophes should be replaced with the right single quote.
			from: "John O'Brian",
			to:   "John O’Brian",
		},
	}

	for _, test := range testCases {
		to := decodeAuthors(test.from)
		if to != test.to {
			t.Errorf("Expected\n%s\ngot\n%s", test.to, to)
		}
	}
}
