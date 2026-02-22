package lexer_test

import (
	"errors"
	"testing"

	lexer "termrome.io/lexer"
	html_lexer "termrome.io/lexer/html"
)

func TestStateMachine(t *testing.T) {
	stateMachine := lexer.NewStateMachine()

	fromState := html_lexer.SEP
	toState := html_lexer.ANCHOR_START
	expectedStateStartIdx := 10
	expectedContextDataKey := "test"
	expectedContextDataValue := "hello world"
	expectedQuoteChar := '"'
	expectedWithinQuotes := false

	stateMachine.Transition(lexer.NewTransitionOption(html_lexer.SEP, 0))

	stateMachine.SetContextData(expectedContextDataKey, expectedContextDataValue)
	stateMachine.SetQuotesState(expectedQuoteChar, expectedWithinQuotes)

	methodCalledCounter := 0

	stateMachine.OnTransition = func(previousState, newState lexer.Token, stateStartIdx int) error {
		if previousState.GetName() != fromState.GetName() {
			t.Errorf("Expected previousState to be %s, received %s instead", fromState.GetName(), previousState.GetName())
		}

		if newState.GetName() != toState.GetName() {
			t.Errorf("Expected newState to be %s, received %s instead", fromState.GetName(), newState.GetName())
		}

		if stateStartIdx != expectedStateStartIdx {
			t.Errorf("Expected stateStartIdx to be %d, received %d instead", expectedStateStartIdx, stateStartIdx)
		}
		methodCalledCounter += 1

		return nil
	}

	stateMachine.Transition(lexer.NewTransitionOption(toState, expectedStateStartIdx))

	if methodCalledCounter != 1 {
		t.Errorf("Expected onTransition to be called when Transition is called")
	}

	currState := stateMachine.GetState()

	if currState.GetName() != toState.GetName() {
		t.Errorf("Expected current start after transition to be %s, received %s instead", toState.GetName(), currState.GetName())
	}

	if stateMachine.GetStateStartIdx() != expectedStateStartIdx {
		t.Errorf("Expected stateStartIdx after transition to be %d, received %d instead", expectedStateStartIdx, stateMachine.GetStateStartIdx())
	}

	value, exists := stateMachine.GetContextData(expectedContextDataKey)
	quotesState := stateMachine.GetQuoteState()

	if !exists {
		t.Errorf("Expected %s to be present in contextData after transition", expectedContextDataKey)
	}

	if value != expectedContextDataValue {
		t.Errorf("Expected contextData[%s] to be %s, recevied %s instead", expectedContextDataKey, expectedContextDataValue, value)
	}

	if *quotesState.GetQuoteChar() != expectedQuoteChar {
		t.Errorf("Expected quoteChar to be %v after transition, received %v instead", expectedQuoteChar, *quotesState.GetQuoteChar())
	}

	if quotesState.GetWithinQuotes() != expectedWithinQuotes {
		t.Errorf("Expected withinQuotes to be %v after transition, received %v instead", expectedWithinQuotes, quotesState.GetWithinQuotes())
	}

	// try with silent transition, shouldn't increase call count
	stateMachine.Transition(lexer.NewTransitionOptionSilent(html_lexer.ANCHOR_END, 20))

	currState = stateMachine.GetState()
	toState = html_lexer.ANCHOR_END
	expectedStateStartIdx = 20

	if methodCalledCounter != 1 {
		t.Errorf("Expected OnTransition to not have been called on silent transition")
	}

	if currState.GetName() != toState.GetName() {
		t.Errorf("Expected current start after transition to be %s, received %s instead", toState.GetName(), currState.GetName())
	}

	if stateMachine.GetStateStartIdx() != expectedStateStartIdx {
		t.Errorf("Expected stateStartIdx after transition to be %d, received %d instead", expectedStateStartIdx, stateMachine.GetStateStartIdx())
	}
}

func TestContextData(t *testing.T) {
	stateMachine := lexer.NewStateMachine()

	_, exists := stateMachine.GetContextData("test")

	if exists {
		t.Errorf("Expected contextData to be empty")
	}

	expectedKey := "test"
	expectedValue := "hello world"

	stateMachine.SetContextData(expectedKey, expectedValue)

	value, exists := stateMachine.GetContextData(expectedKey)

	if !exists {
		t.Errorf("Expected %s key to exists in contextData", expectedKey)
	}

	if value != expectedValue {
		t.Errorf("Expected contextData[%s] to have %s as value, received %s instead", expectedKey, expectedValue, value)
	}
}

func TestQuotesState(t *testing.T) {
	stateMachine := lexer.NewStateMachine()

	quotesState := stateMachine.GetQuoteState()

	expectedQuoteChar := '"'
	expectedWithinQuotes := true

	if quotesState.GetQuoteChar() != nil {
		t.Errorf("Expected quoteChar to be null on an empty state machine")
	}

	if quotesState.GetWithinQuotes() != false {
		t.Errorf("Expected withinQuotes to be false on an empty state machine")
	}

	stateMachine.SetQuotesState(expectedQuoteChar, expectedWithinQuotes)

	quotesState = stateMachine.GetQuoteState()

	if *quotesState.GetQuoteChar() != expectedQuoteChar {
		t.Errorf("Expected quoteChar to be %v, received %v instead", expectedQuoteChar, *quotesState.GetQuoteChar())
	}

	if quotesState.GetWithinQuotes() != expectedWithinQuotes {
		t.Errorf("Expected withinQuotes to be %v, received %v instead", expectedWithinQuotes, quotesState.GetWithinQuotes())
	}
}

func TestOnTrasitionError(t *testing.T) {
	stateMachine := lexer.NewStateMachine()

	stateMachine.OnTransition = func(t1, t2 lexer.Token, i int) error {
		return errors.New("test error")
	}

	err := stateMachine.Transition(lexer.NewTransitionOption(html_lexer.ANCHOR_END, 0))

	if err == nil {
		t.Errorf("Expected OnTransition error to bubble up via Transition")
	}
}
