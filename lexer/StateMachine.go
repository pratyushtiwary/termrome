package lexer

import "sync"

// structs
type QuotesState struct {
	quoteChar    *rune
	withinQuotes bool
}

type StateMachine struct {
	state         Token
	stateStartIdx int
	contextData   map[string]string
	quotesState   *QuotesState
	OnTransition  func(Token, Token, int) error
	rwLock        sync.RWMutex
}

type TransitionOptions struct {
	newState      Token
	stateStartIdx int
	silent        bool
}

// methods

func NewTransitionOptions(newState Token, stateStartIdx int, silent bool) TransitionOptions {
	return TransitionOptions{
		newState:      newState,
		stateStartIdx: stateStartIdx,
		silent:        silent,
	}
}

// QuotesState
func (quotesState *QuotesState) GetQuoteChar() *rune {
	return quotesState.quoteChar
}

func (quotesState *QuotesState) GetWithinQuotes() bool {
	return quotesState.withinQuotes
}

// StateMachine
func (stateMachine *StateMachine) Transition(options TransitionOptions) error {
	if !options.silent && stateMachine.OnTransition != nil {
		err := stateMachine.OnTransition(stateMachine.state, options.newState, options.stateStartIdx)

		if err != nil {
			return err
		}
	}
	stateMachine.state = options.newState
	stateMachine.stateStartIdx = options.stateStartIdx
	return nil
}

func (stateMachine *StateMachine) GetState() Token {
	return stateMachine.state
}

func (stateMachine *StateMachine) GetStateStartIdx() int {
	return stateMachine.stateStartIdx
}

func (stateMachine *StateMachine) GetContextData(key string) (string, bool) {
	stateMachine.rwLock.RLock()
	value, exists := stateMachine.contextData[key]
	stateMachine.rwLock.RUnlock()
	return value, exists
}

func (stateMachine *StateMachine) GetQuoteState() *QuotesState {
	return stateMachine.quotesState
}

func (stateMachine *StateMachine) SetContextData(key string, value string) {
	stateMachine.rwLock.Lock()
	defer stateMachine.rwLock.Unlock()
	stateMachine.contextData[key] = value
}

func (stateMachine *StateMachine) SetQuotesState(quoteChar rune, withinQuotes bool) {
	stateMachine.quotesState = &QuotesState{
		quoteChar:    &quoteChar,
		withinQuotes: withinQuotes,
	}
}

func NewStateMachine() StateMachine {
	return StateMachine{
		contextData: make(map[string]string),
		quotesState: &QuotesState{
			withinQuotes: false,
		},
		rwLock: sync.RWMutex{},
	}
}

func NewTransitionOption(newState Token, stateStartIdx int) TransitionOptions {
	return TransitionOptions{
		newState:      newState,
		stateStartIdx: stateStartIdx,
	}
}

func NewTransitionOptionSilent(newState Token, stateStartIdx int) TransitionOptions {
	transitionOption := NewTransitionOption(newState, stateStartIdx)

	transitionOption.silent = true

	return transitionOption
}
