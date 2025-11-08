package lexer_test

import (
	"testing"

	"termrome.io/lexer"
)

func TestToken(t *testing.T) {
	testToken := lexer.NewToken("test", 't')

	if testToken.GetName() != "test" {
		t.Errorf("Expected test token's name to be 'test', recevied %s instead", testToken.GetName())
	}

	if testToken.GetChar() != 't' {
		t.Errorf("Expected test token's char to be 't', recevied %v instead", testToken.GetChar())
	}
}

func TestTokenWithoutChar(t *testing.T) {
	testToken := lexer.NewTokenWithoutChar("test")

	if testToken.GetName() != "test" {
		t.Errorf("Expected test token's name to be 'test', recevied %s instead", testToken.GetName())
	}

	if testToken.GetChar() != rune(0) {
		t.Errorf("Expected test token's char to be \\0, recevied %v instead", testToken.GetChar())
	}
}
