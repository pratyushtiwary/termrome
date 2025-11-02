package lexer

var ANCHOR_START = NewToken("ANCHOR_START", '<')
var ANCHOR_END = NewToken("ANCHOR_END", '>')
var END_CHAR = NewToken("END_CHAR", '/')
var ATTR_SEP = NewToken("ATTR_SEP", '=')
var SEP = NewToken("SEP", ' ')
var DOUBLE_QUOTE = NewToken("DOUBLE_QUOTE", '"')
var SINGLE_QUOTE = NewToken("SINGLE_QUOTE", '\'')
var TAG = NewTokenWithoutChar("TAG")
var ATTR = NewTokenWithoutChar("ATTR")
var VALUE = NewTokenWithoutChar("VALUE")
var CONTENT = NewTokenWithoutChar("CONTENT")

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
