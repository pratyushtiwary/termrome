package html_lexer_test

import (
	"testing"

	html_lexer "termrome.io/lexer/html"
)

func TestGetValue(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		startIdx  int
		endIdx    int
		quoteChar *rune
		output    string
	}{
		{
			name:      "basic substring",
			input:     "This is a test input",
			startIdx:  10,
			endIdx:    20,
			quoteChar: nil,
			output:    "test input",
		},
		{
			name:      "single quotes case 1",
			input:     "This is a 'test input'",
			startIdx:  10,
			endIdx:    22,
			quoteChar: nil,
			output:    "test input",
		},
		{
			name:      "single quotes case 2",
			input:     "This is a 'test input'",
			startIdx:  10,
			endIdx:    21,
			quoteChar: nil,
			output:    "test input",
		},
		{
			name:      "nested single quote",
			input:     "This is a 'test 'input'",
			startIdx:  10,
			endIdx:    22,
			quoteChar: nil,
			output:    "test 'input",
		},
		{
			name:      "double quotes with escaped single quote",
			input:     `This is a "test 'input"`,
			startIdx:  10,
			endIdx:    22,
			quoteChar: nil,
			output:    "test 'input",
		},
		{
			name:      "with single quote char specified",
			input:     `This is a "test 'input"`,
			startIdx:  10,
			endIdx:    23,
			quoteChar: runePtr('\''),
			output:    `"test 'input"`,
		},
		{
			name:      "with double quote char specified",
			input:     `This is a 'test "input'`,
			startIdx:  10,
			endIdx:    23,
			quoteChar: runePtr('"'),
			output:    `'test "input'`,
		},
		{
			name:      "with spaces and double quote char",
			input:     `This is a  'test "input'  `,
			startIdx:  10,
			endIdx:    26,
			quoteChar: runePtr('"'),
			output:    `'test "input'`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := html_lexer.GetValue(tt.input, tt.startIdx, tt.endIdx, tt.quoteChar)
			if result != tt.output {
				t.Errorf("GetValue() = %q, want %q", result, tt.output)
			}
		})
	}
}

func TestGetAttr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		startIdx int
		endIdx   int
		output   string
	}{
		{
			name:     "basic substring",
			input:    "This is a test input",
			startIdx: 10,
			endIdx:   20,
			output:   "test input",
		},
		{
			name:     "substring with surrounding spaces",
			input:    "This is a  test input  ",
			startIdx: 10,
			endIdx:   23,
			output:   "test input",
		},
		{
			name:     "substring with quotes",
			input:    "This is a 'test input'",
			startIdx: 10,
			endIdx:   22,
			output:   "'test input'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := html_lexer.GetAttr(tt.input, tt.startIdx, tt.endIdx)
			if result != tt.output {
				t.Errorf("GetAttr() = %q, want %q", result, tt.output)
			}
		})
	}
}

func TestRemoveQuotes(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		quoteChar *rune
		output    string
	}{
		{
			name:      "double quotes - matching",
			input:     `"Hello world"`,
			quoteChar: runePtr('"'),
			output:    "Hello world",
		},
		{
			name:      "single quotes - matching",
			input:     "'Hello world'",
			quoteChar: runePtr('\''),
			output:    "Hello world",
		},
		{
			name:      "single quote in middle",
			input:     "'Hello 'world",
			quoteChar: runePtr('\''),
			output:    "Hello 'world",
		},
		{
			name:      "double quote in middle",
			input:     `"Hello "world`,
			quoteChar: runePtr('"'),
			output:    `Hello "world`,
		},
		{
			name:      "single quotes not at start",
			input:     "Hello 'world'",
			quoteChar: runePtr('\''),
			output:    "Hello 'world",
		},
		{
			name:      "double quotes not at start",
			input:     `Hello "world"`,
			quoteChar: runePtr('"'),
			output:    `Hello "world`,
		},
		{
			name:      "no quotes with quote char",
			input:     "Hello world",
			quoteChar: runePtr('"'),
			output:    "Hello world",
		},
		{
			name:      "single quotes - no quote char specified",
			input:     "'Hello world'",
			quoteChar: nil,
			output:    "Hello world",
		},
		{
			name:      "double quotes - no quote char specified",
			input:     `"Hello world"`,
			quoteChar: nil,
			output:    "Hello world",
		},
		{
			name:      "double quotes with wrong quote char",
			input:     `"Hello world"`,
			quoteChar: runePtr('\''),
			output:    `"Hello world"`,
		},
		{
			name:      "single quotes with wrong quote char",
			input:     "'Hello world'",
			quoteChar: runePtr('"'),
			output:    "'Hello world'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := html_lexer.RemoveQuotes(tt.input, tt.quoteChar)
			if result != tt.output {
				t.Errorf("RemoveQuotes() = %q, want %q", result, tt.output)
			}
		})
	}
}

// Helper function to create a pointer to a rune
func runePtr(r rune) *rune {
	return &r
}
