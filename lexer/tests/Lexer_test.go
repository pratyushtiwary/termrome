package lexer_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"termrome.io/lexer"
)

func TestBasicParseHtml(t *testing.T) {
	tests := []struct {
		title string
		input string
	}{
		{
			title: "Empty string",
			input: ``,
		},
		{
			title: "Single tag, without attrs",
			input: `<p>Test</p>`,
		},
		{
			title: "Single tag, with attrs",
			input: `<p abc xyz="hello world" test=1>Test</p>`,
		},
		{
			title: "Nested tags, without attrs",
			input: `<p>Test <h1>Hello World</h1></p>`,
		},
		{
			title: "Nested tags, with attrs",
			input: `<p test="hello world" xyz>Test <h1 abc=1 pqr="h1 test">Hello World</h1></p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			result, err := lexer.ParseHtml(tt.input)
			if err != nil {
				t.Errorf("Expected nodes array received %s error instead", err.Error())
				return
			}
			snaps.MatchSnapshot(t, result)
		})
	}
}

func TestParseHtmlAdvanced(t *testing.T) {
	tests := []struct {
		title string
		input string
	}{
		{
			title: "void tag, not nested",
			input: `<img src="https://example.com" disabled width=200>`,
		},
		{
			title: "void tag, nested",
			input: `<p><img src="https://example.com" disabled width=200> test</p>`,
		},
		{
			title: "user-defined void tag, not nested",
			input: `<xyz abc=1 pqr="hello world" test> test`,
		},
		{
			title: "user-defined void tag, nested",
			input: `<p><xyz abc=1 pqr="hello world" test> <abc>test</p>`,
		},
		{
			title: "without root",
			input: `<b>Hello world</b>test <p>Hello world</p>`,
		},
		{
			title: "dangling <",
			input: `<p>Text < <b>bold</b></p>`,
		},
		{
			title: "simple html file",
			input: `<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Test Simple</title>
</head>
<body>
	<h1>Hello world</h1>
</body>
</html>`,
		},
		{
			title: "highly nested html file",
			input: `<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Test highly nested</title>
</head>
<body>
	<h1>Hello world <p>I'm a <b>p <i><xyz>tag</i></b></p></h1>
</body>
</html>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			result, err := lexer.ParseHtml(tt.input)
			if err != nil {
				t.Fatalf("ParseHtml() error = %v", err)
			}
			snaps.MatchSnapshot(t, result)
		})
	}
}
