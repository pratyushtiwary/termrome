package html_lexer_test

import (
	"testing"

	html_lexer "termrome.io/lexer/html"
)

func TestTextNode(t *testing.T) {
	textNode := html_lexer.NewTextNode("Hello\r\nWorld")

	if textNode.Content != "Hello\nWorld" {
		t.Errorf("Expected textNode's content to be Hello\nWorld, received %s instead", textNode.Content)
	}

	if textNode.Type != textNode.GetType() || textNode.GetType() != html_lexer.TEXT_NODE {
		t.Errorf("Expected textNode's type to be %s, received %s instead", html_lexer.TEXT_NODE, textNode.Type)
	}

	if textNode.HasEnding != false {
		t.Errorf("Expected textNode's hasEnding to be false, received %v instead", textNode.HasEnding)
	}
}

func TestElementNode(t *testing.T) {
	elementNode, err := html_lexer.NewElementNode("testTag")

	if err != nil {
		t.Errorf("Expected elementNode to be created succesfully, received an error instead, error: %s", err.Error())
	}

	if elementNode.Tag != "testTag" {
		t.Errorf("Expected elementNode's tag to be 'testTag', received %s instead", elementNode.Tag)
	}

	if elementNode.HasEnding != false {
		t.Errorf("Expected elementNode's hasEnding to be false, received %v instead", elementNode.HasEnding)
	}

	if elementNode.IsKnownVoid != false {
		t.Errorf("Expected elementNode's IsKnownVoid to be false, received %v instead", elementNode.IsKnownVoid)
	}

	if elementNode.Type != elementNode.GetType() || elementNode.GetType() != html_lexer.ELEMENT_NODE {
		t.Errorf("Expected elementNode's type to be %s, received %s instead", html_lexer.ELEMENT_NODE, elementNode.Type)
	}
}

func TestKnownVoidTag(t *testing.T) {
	elementNode, err := html_lexer.NewElementNode("br")

	if err != nil {
		t.Errorf("Expected elementNode to be created succesfully, received an error instead, error: %s", err.Error())
	}

	if elementNode.Tag != "br" {
		t.Errorf("Expected elementNode's tag to be 'br', received %s instead", elementNode.Tag)
	}

	if elementNode.HasEnding != false {
		t.Errorf("Expected elementNode's hasEnding to be false, received %v instead", elementNode.HasEnding)
	}

	if elementNode.IsKnownVoid != true {
		t.Errorf("Expected elementNode's IsKnownVoid to be true, received %v instead", elementNode.IsKnownVoid)
	}

	if elementNode.Type != elementNode.GetType() || elementNode.GetType() != html_lexer.ELEMENT_NODE {
		t.Errorf("Expected elementNode's type to be %s, received %s instead", html_lexer.ELEMENT_NODE, elementNode.Type)
	}
}

func TestElementNodeShouldReturnError(t *testing.T) {
	_, err := html_lexer.NewElementNode("")

	if err == nil {
		t.Error("Expected to received error on element node creation")
	}
}

func TestElementNodeData(t *testing.T) {
	elementNode, err := html_lexer.NewElementNode("testTag")

	if err != nil {
		t.Errorf("Expected elementNode to be created succesfully, received an error instead, error: %s", err.Error())
	}

	_, exists := elementNode.GetData("test")

	if exists {
		t.Error("Expected test to not exist in elementNode's data")
	}

	elementNode.SetData("test", html_lexer.Data{
		Value:     "1",
		IsPresent: true,
	})

	value, exists := elementNode.GetData("test")

	if !exists {
		t.Error("Expected test to exist in elementNode's data")
	}

	if value.Value != "1" {
		t.Errorf("Expected value for test to be 1 received %s instead", value.Value)
	}

	if !value.IsPresent {
		t.Errorf("Expected IsPresent for test to be true received %v instead", value.IsPresent)
	}

}
