package lexer

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
	onTransition  func(Token, Token, int)
}

type TransitionOptions struct {
	newState      Token
	stateStartIdx int
	silent        bool
}

// methods

// QuotesState
func (quotesState *QuotesState) GetQuoteChar() *rune {
	return quotesState.quoteChar
}

func (quotesState *QuotesState) GetWithinQuotes() bool {
	return quotesState.withinQuotes
}

// StateMachine
func (stateMachine *StateMachine) Transition(options TransitionOptions) {
	if !options.silent && stateMachine.onTransition != nil {
		stateMachine.onTransition(stateMachine.state, options.newState, options.stateStartIdx)
	}
	stateMachine.state = options.newState
	stateMachine.stateStartIdx = options.stateStartIdx
}

func (stateMachine *StateMachine) GetState() Token {
	return stateMachine.state
}

func (stateMachine *StateMachine) GetStateStartIdx() int {
	return stateMachine.stateStartIdx
}

func (stateMachine *StateMachine) GetContext() map[string]string {
	return stateMachine.contextData
}

func (stateMachine *StateMachine) GetQuoteState() *QuotesState {
	return stateMachine.quotesState
}

func (stateMachine *StateMachine) SetContextData(key string, value string) {
	stateMachine.contextData[key] = value
}

func (stateMachine *StateMachine) SetQuotesState(quoteChar rune, withinQuotes bool) {
	stateMachine.quotesState = &QuotesState{
		quoteChar:    &quoteChar,
		withinQuotes: withinQuotes,
	}
}
