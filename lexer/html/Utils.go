package html

import (
	"strings"

	lexer "termrome.io/lexer"
)

// recovery checkers
var RequiresTagStartRecovery = lexer.NewTransitionStrategyBuilder().
	StateEq(&ANCHOR_START).
	CurrCharEq(SEP.GetChar(), END_CHAR.GetChar(), ANCHOR_END.GetChar(), ANCHOR_START.GetChar()).
	Build()

var RequiresAttrRecovery = lexer.NewTransitionStrategyBuilder().
	StateEq(&ATTR_SEP).
	CurrCharEq(SEP.GetChar()).
	Build()

// normal checkers
var IsAnchorStart = lexer.NewTransitionStrategyBuilder().
	CurrCharEq(ANCHOR_START.GetChar()).
	Build()

var IsAnchorEndChar = lexer.NewTransitionStrategyBuilder().
	StateEq(&ANCHOR_START).
	CurrCharEq(END_CHAR.GetChar()).
	Build()

var IsTagNameStart = lexer.NewTransitionStrategyBuilder().
	StateEq(&ANCHOR_START).
	CurrCharNotEq(SEP.GetChar(), END_CHAR.GetChar(), ANCHOR_END.GetChar(), ANCHOR_START.GetChar()).
	Build()

var IsAnchorTagEnd = lexer.NewTransitionStrategyBuilder().
	StateEq(&END_CHAR).
	CurrCharEq(ANCHOR_END.GetChar()).
	Build()

var IsAttrStart = lexer.NewTransitionStrategyBuilder().
	StateEq(&TAG).
	CurrCharEq(SEP.GetChar(), ANCHOR_END.GetChar()).
	Build()

var IsValueStart = lexer.NewTransitionStrategyBuilder().
	StateEq(&ATTR).
	CurrCharEq(
		ANCHOR_END.GetChar(),
		ATTR_SEP.GetChar(),
		SEP.GetChar(),
		DOUBLE_QUOTE.GetChar(),
		SINGLE_QUOTE.GetChar(),
	).
	Build()

var IsAttrSep = lexer.NewTransitionStrategyBuilder().
	StateEq(&ATTR_SEP).
	Build()

var IsValueEnd = lexer.NewTransitionStrategyBuilder().
	StateEq(&VALUE).
	CurrCharEq(SEP.GetChar(), ANCHOR_END.GetChar(), DOUBLE_QUOTE.GetChar(), SINGLE_QUOTE.GetChar()).
	Build()

func RemoveQuotes(inputString string, quoteChar *rune) string {
	if quoteChar != nil {
		inputString = strings.TrimPrefix(inputString, string(*quoteChar))
		inputString = strings.TrimSuffix(inputString, string(*quoteChar))
		return inputString
	}

	doubleQuote := string(DOUBLE_QUOTE.GetChar())
	singleQuote := string(SINGLE_QUOTE.GetChar())

	if strings.HasPrefix(inputString, doubleQuote) || strings.HasPrefix(inputString, singleQuote) {
		inputString = inputString[1:]
	}

	if len(inputString) > 0 && strings.HasSuffix(inputString, doubleQuote) || strings.HasSuffix(inputString, singleQuote) {
		inputString = inputString[:len(inputString)-1]
	}

	return inputString
}

func GetValue(inputString string, startIdx int, endIdx int, quoteChar *rune) string {
	return strings.TrimSpace(RemoveQuotes(inputString[startIdx:endIdx], quoteChar))
}

func GetAttr(inputString string, startIdx int, endIdx int) string {
	return strings.TrimSpace(inputString[startIdx:endIdx])
}

func IsVoidTag(tagName string) bool {
	if _, exists := VOID_TAGS[tagName]; exists {
		return true
	}

	return false
}
