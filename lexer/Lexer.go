package lexer

import (
	"fmt"
)

func handleRecovery(stateMachine *StateMachine) error {
	stateStartIdx := stateMachine.GetStateStartIdx()
	return stateMachine.Transition(TransitionOptions{
		newState:      CONTENT,
		stateStartIdx: stateStartIdx - 1, // -1 to take < into account for content as well
		silent:        true,
	})
}

func ParseHtml(inputString string) ([]Node, error) {
	var err error

	tagStack := make([]*ElementNode, 0)
	nodes := make([]Node, 0)
	stateMachine := NewStateMachine()

	popTagStack := func() *ElementNode {
		var node *ElementNode
		if len(tagStack) > 0 {
			node = tagStack[len(tagStack)-1]
			tagStack = tagStack[:len(tagStack)-1]
		}

		return node
	}

	stateMachine.Transition(TransitionOptions{
		newState:      CONTENT,
		stateStartIdx: 0,
		silent:        true,
	})

	stateMachine.OnTransition = func(oldState, newState Token, stateStartIdx int) error {
		if newState == CONTENT && inputString[stateStartIdx-1] == byte(ANCHOR_END.GetChar()) && len(tagStack) > 0 {
			currNode := tagStack[len(tagStack)-1]

			if currNode == nil {
				return nil
			}

			if currNode.IsKnownVoid {
				node := popTagStack()
				if len(tagStack) > 0 {
					tagStack[len(tagStack)-1].Children = append(tagStack[len(tagStack)-1].Children, node)
				} else {
					nodes = append(nodes, node)
				}
			}
		}

		if oldState == CONTENT {
			content := inputString[stateMachine.GetStateStartIdx() : stateStartIdx-1]

			if content == "" {
				return nil
			}

			textNode := NewTextNode(content)

			if len(tagStack) == 0 {
				nodes = append(nodes, textNode)
				return nil
			}

			tagStack[len(tagStack)-1].Children = append(tagStack[len(tagStack)-1].Children, textNode)
		}

		return nil
	}

	for currCharIdx, currChar := range inputString {

		// anchor start logic
		if IsAnchorStart(&stateMachine, currChar) {
			err = stateMachine.Transition(TransitionOptions{
				newState:      ANCHOR_START,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})

			if err != nil {
				return nil, err
			}
			continue
		}

		if IsAnchorEndChar(&stateMachine, currChar) {
			err = stateMachine.Transition(TransitionOptions{
				newState:      END_CHAR,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})

			if err != nil {
				return nil, err
			}
			continue
		}

		if IsTagNameStart(&stateMachine, currChar) {
			err = stateMachine.Transition(TransitionOptions{
				newState:      TAG,
				stateStartIdx: currCharIdx,
				silent:        false,
			})

			if err != nil {
				return nil, err
			}
			continue
		}

		// recovery logic
		if RequiresTagStartRecovery(&stateMachine, currChar) {
			err = handleRecovery(&stateMachine)

			if err != nil {
				return nil, err
			}

			err = stateMachine.Transition(TransitionOptions{
				newState:      CONTENT,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})

			if err != nil {
				return nil, err
			}
			continue
		}

		if RequiresAttrRecovery(&stateMachine, currChar) {
			continue
		}

		// attr and tag end logic
		if IsAnchorTagEnd(&stateMachine, currChar) {
			endingTagName := inputString[stateMachine.GetStateStartIdx():currCharIdx]
			node := popTagStack()
			var prevNode *ElementNode

			if node != nil {
				node.HasEnding = node.Tag == endingTagName
			}

			for len(tagStack) > 0 && node.Tag != endingTagName {
				prevNode = node
				node = popTagStack()
				node.HasEnding = node.Tag == endingTagName
				node.Children = append(node.Children, prevNode)
			}

			if len(tagStack) > 0 {
				tagStack[len(tagStack)-1].Children = append(tagStack[len(tagStack)-1].Children, node)
			} else if node != nil {
				nodes = append(nodes, node)
			}
			err = stateMachine.Transition(TransitionOptions{
				newState:      CONTENT,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})

			if err != nil {
				return nil, err
			}
			continue
		}

		if IsAttrStart(&stateMachine, currChar) {
			var elementNode *ElementNode
			tagName := inputString[stateMachine.GetStateStartIdx():currCharIdx]
			elementNode, err = NewElementNode(tagName)

			if err != nil {
				return nil, err
			}

			tagStack = append(tagStack, elementNode)
			var nextState = ATTR

			if currChar == ANCHOR_END.GetChar() {
				nextState = CONTENT
			}

			err = stateMachine.Transition(TransitionOptions{
				newState:      nextState,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})

			if err != nil {
				return nil, err
			}
			continue
		}

		if IsValueStart(&stateMachine, currChar) {
			attrName := GetAttr(
				inputString,
				stateMachine.GetStateStartIdx(),
				currCharIdx,
			)

			if currChar == ANCHOR_END.GetChar() && attrName == "" {
				err = stateMachine.Transition(TransitionOptions{
					newState:      CONTENT,
					stateStartIdx: currCharIdx + 1,
					silent:        false,
				})

				if err != nil {
					return nil, err
				}
				continue
			}

			if attrName == "" {
				continue
			}

			node := tagStack[len(tagStack)-1]
			node.SetData(attrName, Data{
				IsPresent: true,
				Value:     "",
			})

			var nextState = ATTR_SEP

			if currChar == ANCHOR_END.GetChar() {
				nextState = CONTENT
			} else if currChar == SEP.GetChar() {
				nextState = ATTR
			}

			err = stateMachine.Transition(TransitionOptions{
				newState:      nextState,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})

			if err != nil {
				return nil, err
			}
			stateMachine.SetContextData("attrName", attrName)
			continue
		}

		if IsAttrSep(&stateMachine, currChar) {
			quoteChar := currChar
			caretIdx := currCharIdx

			if currChar == ANCHOR_END.GetChar() {
				err = stateMachine.Transition(TransitionOptions{
					newState:      CONTENT,
					stateStartIdx: currCharIdx + 1,
					silent:        false,
				})

				if err != nil {
					return nil, err
				}
				continue
			}

			if quoteChar == DOUBLE_QUOTE.GetChar() || quoteChar == SINGLE_QUOTE.GetChar() {
				stateMachine.SetQuotesState(quoteChar, true)
				caretIdx += 1
			}
			err = stateMachine.Transition(TransitionOptions{
				newState:      VALUE,
				stateStartIdx: caretIdx,
				silent:        false,
			})

			if err != nil {
				return nil, err
			}
			continue
		}

		if IsValueEnd(&stateMachine, currChar) {
			quotesState := stateMachine.GetQuoteState()

			if quotesState != nil && quotesState.GetWithinQuotes() && currChar != *quotesState.GetQuoteChar() {
				continue
			}

			var quoteChar *rune

			if quotesState != nil {
				quoteChar = quotesState.GetQuoteChar()
			}

			value := GetValue(
				inputString,
				stateMachine.GetStateStartIdx(),
				currCharIdx,
				quoteChar,
			)

			attrName, attrExists := stateMachine.GetContextData("attrName")

			if !attrExists {
				return nil, fmt.Errorf("failed to fetch attribute name for value: %s", value)
			}

			node := tagStack[len(tagStack)-1]
			data, exists := node.GetData(attrName)
			if exists {
				data.Value = value
				node.SetData(attrName, data)
			} else {
				node.SetData(attrName, Data{
					Value:     value,
					IsPresent: true,
				})
			}

			stateMachine.SetQuotesState(' ', false)
			var nextState = ATTR

			if currChar == ANCHOR_END.GetChar() {
				nextState = CONTENT
			}

			err = stateMachine.Transition(TransitionOptions{
				newState:      nextState,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})

			if err != nil {
				return nil, err
			}
		}
	}

	if len(tagStack) > 0 {
		node := popTagStack()
		var prevNode *ElementNode

		for len(tagStack) > 0 {
			prevNode = node
			node = popTagStack()
			node.Children = append(node.Children, prevNode)
		}

		nodes = append(nodes, node)
	}

	remainingContent := inputString[stateMachine.GetStateStartIdx():]
	if remainingContent != "" {
		lastNode := nodes[len(nodes)-1]
		textNode := NewTextNode(remainingContent)
		if lastNode != nil {
			if elemNode, ok := lastNode.(*ElementNode); ok {
				if !elemNode.HasEnding {
					elemNode.Children = append(elemNode.Children, textNode)
				}
			} else {
				nodes = append(nodes, textNode)
			}
		}
	}

	return nodes, nil
}
