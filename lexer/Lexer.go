package lexer

func handleRecovery(stateMachine StateMachine) {
	stateStartIdx := stateMachine.GetStateStartIdx()
	stateMachine.Transition(TransitionOptions{
		newState:      CONTENT,
		stateStartIdx: stateStartIdx,
		silent:        true,
	})
}

func ParseHtml(inputString string) []Node {
	tagStack := make([]*ElementNode, 0)
	nodes := make([]Node, 0)
	stateMachine := StateMachine{
		contextData: map[string]string{},
	}

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

	stateMachine.onTransition = func(oldState, newState Token, stateStartIdx int) {
		if newState == CONTENT && inputString[stateStartIdx-1] == byte(ANCHOR_END.GetChar()) && len(tagStack) > 0 {
			currNode := tagStack[len(tagStack)-1]

			if currNode == nil {
				return
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
				return
			}

			if len(tagStack) == 0 {
				nodes = append(nodes, NewTextNode(content))
				return
			}
			tagStack[len(tagStack)-1].Children = append(tagStack[len(tagStack)-1].Children, NewTextNode(content))
		}
	}

	for currCharIdx, currChar := range inputString {

		// anchor start logic
		if IsAnchorStart(stateMachine, currChar) {
			stateMachine.Transition(TransitionOptions{
				newState:      ANCHOR_START,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})
			continue
		}

		if IsAnchorEndChar(stateMachine, currChar) {
			stateMachine.Transition(TransitionOptions{
				newState:      END_CHAR,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})
			continue
		}

		if IsTagNameStart(stateMachine, currChar) {
			stateMachine.Transition(TransitionOptions{
				newState:      TAG,
				stateStartIdx: currCharIdx,
				silent:        false,
			})
			continue
		}

		// recovery logic
		if RequiresTagStartRecovery(stateMachine, currChar) {
			handleRecovery(stateMachine)
			stateMachine.Transition(TransitionOptions{
				newState:      CONTENT,
				stateStartIdx: currCharIdx,
				silent:        false,
			})
			continue
		}

		if RequiresAttrRecovery(stateMachine, currChar) {
			continue
		}

		// attr and tag end logic
		if IsAnchorTagEnd(stateMachine, currChar) {
			endingTagName := inputString[stateMachine.GetStateStartIdx():currCharIdx]
			node := popTagStack()
			var prevNode *ElementNode = nil

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
			stateMachine.Transition(TransitionOptions{
				newState:      CONTENT,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})
			continue
		}

		if IsAttrStart(stateMachine, currChar) {
			tagName := inputString[stateMachine.GetStateStartIdx():currCharIdx]
			tagStack = append(tagStack, NewElementNode(tagName))
			var nextState Token = ATTR

			if currChar == ANCHOR_END.GetChar() {
				nextState = CONTENT
			}

			stateMachine.Transition(TransitionOptions{
				newState:      nextState,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})
			continue
		}

		if IsValueStart(stateMachine, currChar) {
			attrName := GetAttr(
				inputString,
				stateMachine.GetStateStartIdx(),
				currCharIdx,
			)

			if currChar == ANCHOR_END.GetChar() && attrName == "" {
				stateMachine.Transition(TransitionOptions{
					newState:      CONTENT,
					stateStartIdx: currCharIdx + 1,
					silent:        false,
				})
				continue
			}

			if attrName == "" {
				continue
			}

			node := tagStack[len(tagStack)-1]
			node.Data[attrName] = Data{
				IsPresent: true,
				Value:     "",
			}

			var nextState Token = ATTR_SEP

			if currChar == ANCHOR_END.GetChar() {
				nextState = CONTENT
			} else if currChar == SEP.GetChar() {
				nextState = ATTR
			}

			stateMachine.Transition(TransitionOptions{
				newState:      nextState,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})
			stateMachine.SetContextData("attrName", attrName)
			continue
		}

		if IsAttrSep(stateMachine, currChar) {
			quoteChar := currChar
			caretIdx := currCharIdx

			if currChar == ANCHOR_END.GetChar() {
				stateMachine.Transition(TransitionOptions{
					newState:      CONTENT,
					stateStartIdx: currCharIdx + 1,
					silent:        false,
				})
				continue
			}

			if quoteChar == DOUBLE_QUOTE.GetChar() || quoteChar == SINGLE_QUOTE.GetChar() {
				stateMachine.SetQuotesState(quoteChar, true)
				caretIdx += 1
			}
			stateMachine.Transition(TransitionOptions{
				newState:      VALUE,
				stateStartIdx: caretIdx,
				silent:        false,
			})
			continue
		}

		if IsValueEnd(stateMachine, currChar) {
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

			attrName := stateMachine.GetContext()["attrName"]
			node := tagStack[len(tagStack)-1]
			if data, exists := node.Data[attrName]; exists {
				data.Value = value
				node.Data[attrName] = data
			} else {
				node.Data[attrName] = Data{
					Value:     value,
					IsPresent: true,
				}
			}

			stateMachine.SetQuotesState(' ', false)
			var nextState Token = ATTR

			if currChar == ANCHOR_END.GetChar() {
				nextState = CONTENT
			}

			stateMachine.Transition(TransitionOptions{
				newState:      nextState,
				stateStartIdx: currCharIdx + 1,
				silent:        false,
			})
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
		if lastNode != nil {
			if elemNode, ok := lastNode.(*ElementNode); ok {
				if !elemNode.HasEnding {
					elemNode.Children = append(elemNode.Children, NewTextNode(remainingContent))
				}
			} else {
				nodes = append(nodes, NewTextNode(remainingContent))
			}
		}
	}

	return nodes
}
