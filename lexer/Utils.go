package lexer

import "strings"

// recovery checkers
var RequiresTagStartRecovery = NewTransitionStrategyBuilder().
	StateEq(&ANCHOR_START).
	CurrCharEq(SEP.char, END_CHAR.char, ANCHOR_END.char, ANCHOR_START.char).
	Build()

var RequiresAttrRecovery = NewTransitionStrategyBuilder().
	StateEq(&ATTR_SEP).
	CurrCharEq(SEP.char).
	Build()

// normal checkers
var IsAnchorStart = NewTransitionStrategyBuilder().
	CurrCharEq(ANCHOR_START.char).
	Build()

var IsAnchorEndChar = NewTransitionStrategyBuilder().
	StateEq(&ANCHOR_START).
	CurrCharEq(END_CHAR.char).
	Build()

var IsTagNameStart = NewTransitionStrategyBuilder().
	StateEq(&ANCHOR_START).
	CurrCharNotEq(SEP.char, END_CHAR.char, ANCHOR_END.char, ANCHOR_START.char).
	Build()

var IsAnchorTagEnd = NewTransitionStrategyBuilder().
	StateEq(&END_CHAR).
	CurrCharEq(ANCHOR_END.char).
	Build()

var IsAttrStart = NewTransitionStrategyBuilder().
	StateEq(&TAG).
	CurrCharEq(SEP.char, ANCHOR_END.char).
	Build()

var IsValueStart = NewTransitionStrategyBuilder().
	StateEq(&ATTR).
	CurrCharEq(
		ANCHOR_END.char,
		ATTR_SEP.char,
		SEP.char,
		DOUBLE_QUOTE.char,
		SINGLE_QUOTE.char,
	).
	Build()

var IsAttrSep = NewTransitionStrategyBuilder().
	StateEq(&ATTR_SEP).
	Build()

var IsValueEnd = NewTransitionStrategyBuilder().
	StateEq(&VALUE).
	CurrCharEq(SEP.char, ANCHOR_END.char, DOUBLE_QUOTE.char, SINGLE_QUOTE.char).
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
