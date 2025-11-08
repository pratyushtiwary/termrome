package lexer

// struct
type TransitionStrategyBuilder struct {
	expectedState *Token
	charEq        map[rune]struct{}
	charNotEq     map[rune]struct{}
}

// methods
func (transitionStrategyBuilder *TransitionStrategyBuilder) StateEq(state *Token) *TransitionStrategyBuilder {
	transitionStrategyBuilder.expectedState = state
	return transitionStrategyBuilder
}

func (transitionStrategyBuilder *TransitionStrategyBuilder) CurrCharEq(currChar ...rune) *TransitionStrategyBuilder {
	for _, char := range currChar {
		transitionStrategyBuilder.charEq[char] = struct{}{}
	}
	return transitionStrategyBuilder
}

func (transitionStrategyBuilder *TransitionStrategyBuilder) CurrCharNotEq(currChar ...rune) *TransitionStrategyBuilder {
	for _, char := range currChar {
		transitionStrategyBuilder.charNotEq[char] = struct{}{}
	}
	return transitionStrategyBuilder
}

func (transitionStrategyBuilder *TransitionStrategyBuilder) Build() func(*StateMachine, rune) bool {
	return func(state *StateMachine, currChar rune) bool {
		matchingCondition := true

		if len(transitionStrategyBuilder.charEq) > 0 {
			if _, exists := transitionStrategyBuilder.charEq[currChar]; !exists {
				matchingCondition = false
			}
		}

		if len(transitionStrategyBuilder.charNotEq) > 0 {
			if _, exists := transitionStrategyBuilder.charNotEq[currChar]; exists {
				matchingCondition = false
			}
		}

		if transitionStrategyBuilder.expectedState != nil && state.GetState().name != transitionStrategyBuilder.expectedState.name {
			matchingCondition = false
		}

		return matchingCondition
	}
}

func NewTransitionStrategyBuilder() *TransitionStrategyBuilder {
	return &TransitionStrategyBuilder{
		charEq:    make(map[rune]struct{}),
		charNotEq: make(map[rune]struct{}),
	}
}
