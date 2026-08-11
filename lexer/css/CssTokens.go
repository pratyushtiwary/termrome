package css

import (
	lexer "termrome.io/lexer"
)

var ADDRESS_PART = lexer.NewTokenWithoutChar("ADDRESS_PART")
var CLASS_STARTER = lexer.NewToken("CLASS_STARTER", '.')
var ID_STARTER = lexer.NewToken("ID_STARTER", '#')
var KEY_SEPERATOR = lexer.NewToken("KEY_SEPERATOR", ':')
var VALUE_TERMINATOR = lexer.NewToken("VALUE_TERMINATOR", ';')
var BODY_START = lexer.NewToken("BODY_START", '{')
var BODY_CLOSE = lexer.NewToken("BODY_CLOSE", '}')
var SEP = lexer.NewToken("SEP", ' ')
