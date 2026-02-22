package html

import (
	lexer "termrome.io/lexer"
)

var ANCHOR_START = lexer.NewToken("ANCHOR_START", '<')
var ANCHOR_END = lexer.NewToken("ANCHOR_END", '>')
var END_CHAR = lexer.NewToken("END_CHAR", '/')
var ATTR_SEP = lexer.NewToken("ATTR_SEP", '=')
var SEP = lexer.NewToken("SEP", ' ')
var DOUBLE_QUOTE = lexer.NewToken("DOUBLE_QUOTE", '"')
var SINGLE_QUOTE = lexer.NewToken("SINGLE_QUOTE", '\'')
var TAG = lexer.NewTokenWithoutChar("TAG")
var ATTR = lexer.NewTokenWithoutChar("ATTR")
var VALUE = lexer.NewTokenWithoutChar("VALUE")
var CONTENT = lexer.NewTokenWithoutChar("CONTENT")

const TEXT_NODE = "TEXT"
const ELEMENT_NODE = "ELEMENT"

var VOID_TAGS = map[string]struct{}{
	"area":    {},
	"base":    {},
	"br":      {},
	"col":     {},
	"command": {},
	"embed":   {},
	"hr":      {},
	"img":     {},
	"input":   {},
	"keygen":  {},
	"link":    {},
	"meta":    {},
	"param":   {},
	"source":  {},
	"track":   {},
	"wbr":     {},
}
