package main

import (
	"testing"
)

func TestToString(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{
			input:  "Title",
			output: "Title",
		},
		{
			input:  "This is a {Title}",
			output: "This is a Title",
		},
		{
			input:  "This is a {Title}",
			output: "This is a Title",
		},
		{
			input:  `{\#h00t}: Censorship Resistant Microblogging`,
			output: `#h00t: Censorship Resistant Microblogging`,
		},
		{
			input:  "``Good'' Worms and Human Rights",
			output: "“Good” Worms and Human Rights",
		},
		{
			input:  "An Analysis of {China}'s ``{Great Cannon}''",
			output: "An Analysis of China’s “Great Cannon”",
		},
		{
			input:  `lib$\cdot$erate, (n):`,
			output: `lib·erate, (n):`,
		},
		{
			input:  "Well -- Exploring the {Great} {Firewall}'s Poisoned {DNS}",
			output: "Well – Exploring the Great Firewall’s Poisoned DNS",
		},
	}

	for _, test := range testCases {
		output := decodeTitle(test.input)
		if output != test.output {
			t.Errorf("Expected\n%s\ngot\n%s", test.output, output)
		}
	}
}

func TestDecodeAuthors(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{
			// Multiple authors should be separated by commas.
			input:  "John Doe and Jane Doe",
			output: "John Doe, Jane Doe",
		},
		{
			// Single authors should remain as-is.
			input:  "John Doe",
			output: "John Doe",
		},
		{
			// Single-name authors should remain as-is.
			input:  "John and Jane",
			output: "John, Jane",
		},
		{
			// Non-ASCII characters should be unaffected.
			input:  "Jóhn Doe",
			output: "Jóhn Doe",
		},
		{
			// Apostrophes should be replaced with the right single quote.
			input:  "John O'Brian",
			output: "John O’Brian",
		},
	}

	for _, test := range testCases {
		output := decodeAuthors(test.input)
		if output != test.output {
			t.Errorf("Expected\n%s\ngot\n%s", test.output, output)
		}
	}
}
