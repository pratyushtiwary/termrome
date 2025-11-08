package lexer_test

import (
	"testing"

	"termrome.io/lexer"
)

func TestStateEq(t *testing.T) {
	builder := lexer.NewTransitionStrategyBuilder().StateEq(&lexer.ANCHOR_END).Build()
	stateMachine := lexer.NewStateMachine()

	stateMachine.Transition(lexer.NewTransitionOption(lexer.ANCHOR_END, 0))

	if builder(&stateMachine, rune(0)) != true {
		t.Errorf("Expected stateEq(%s) to match with current stateMachine's state", lexer.ANCHOR_END.GetName())
	}

	stateMachine.Transition(lexer.NewTransitionOption(lexer.ANCHOR_START, 0))

	if builder(&stateMachine, rune(0)) != false {
		t.Errorf("Expected stateEq(%s) not to match with current stateMachine's state", lexer.ANCHOR_START.GetName())
	}
}

func TestCharEq(t *testing.T) {
	builder := lexer.NewTransitionStrategyBuilder().StateEq(&lexer.ANCHOR_END).CurrCharEq(lexer.ANCHOR_END.GetChar()).Build()
	stateMachine := lexer.NewStateMachine()

	stateMachine.Transition(lexer.NewTransitionOption(lexer.ANCHOR_END, 0))

	if builder(&stateMachine, lexer.ANCHOR_END.GetChar()) != true {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, %v) to match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_END.GetChar(), currState.GetName(), currState.GetChar())
	}

	if builder(&stateMachine, lexer.ANCHOR_START.GetChar()) != false {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, %v) to not match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_END.GetChar(), currState.GetName(), currState.GetChar())
	}

	stateMachine.Transition(lexer.NewTransitionOption(lexer.ANCHOR_START, 0))

	if builder(&stateMachine, lexer.ANCHOR_START.GetChar()) != false {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, %v) to not match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_END.GetChar(), currState.GetName(), currState.GetChar())
	}

	if builder(&stateMachine, lexer.ANCHOR_END.GetChar()) != false {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, %v) to not match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_END.GetChar(), currState.GetName(), currState.GetChar())
	}
}

func TestCharNotEq(t *testing.T) {
	builder := lexer.NewTransitionStrategyBuilder().StateEq(&lexer.ANCHOR_END).CurrCharNotEq(lexer.ANCHOR_START.GetChar()).Build()
	stateMachine := lexer.NewStateMachine()

	stateMachine.Transition(lexer.NewTransitionOption(lexer.ANCHOR_END, 0))

	if builder(&stateMachine, lexer.ANCHOR_END.GetChar()) != true {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, !%v) to match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_START.GetChar(), currState.GetName(), currState.GetChar())
	}

	if builder(&stateMachine, lexer.ANCHOR_START.GetChar()) != false {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, !%v) to not match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_START.GetChar(), currState.GetName(), currState.GetChar())
	}

	stateMachine.Transition(lexer.NewTransitionOption(lexer.ANCHOR_START, 0))

	if builder(&stateMachine, lexer.ANCHOR_START.GetChar()) != false {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, !%v) to not match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_START.GetChar(), currState.GetName(), currState.GetChar())
	}

	if builder(&stateMachine, lexer.ANCHOR_END.GetChar()) != false {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, !%v) to not match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_START.GetChar(), currState.GetName(), currState.GetChar())
	}
}

func TestWithAllConditions(t *testing.T) {
	builder := lexer.NewTransitionStrategyBuilder().StateEq(&lexer.ANCHOR_END).CurrCharEq(lexer.ANCHOR_START.GetChar()).CurrCharNotEq(lexer.ANCHOR_END.GetChar()).Build()
	stateMachine := lexer.NewStateMachine()

	stateMachine.Transition(lexer.NewTransitionOption(lexer.ANCHOR_END, 0))

	if builder(&stateMachine, lexer.ANCHOR_START.GetChar()) != true {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, %v) to match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_START.GetChar(), currState.GetName(), currState.GetChar())
	}

	if builder(&stateMachine, lexer.ANCHOR_END.GetChar()) != false {
		currState := stateMachine.GetState()
		t.Errorf("Expected (%s, !%v) to match with builder(%s, %v)", lexer.ANCHOR_END.GetName(), lexer.ANCHOR_END.GetChar(), currState.GetName(), currState.GetChar())
	}

}
