package lexer

import (
	"encoding/json"
	"strings"
)

// structs
type Node interface {
	GetType() string
}
type TextNode struct {
	Type      string
	Content   string
	HasEnding bool
}

func (t TextNode) GetType() string { return t.Type }

type Data struct {
	Value     string
	IsPresent bool
}
type ElementNode struct {
	Type        string
	Tag         string
	Children    []Node
	IsKnownVoid bool
	Data        map[string]Data
	HasEnding   bool
}

func (t ElementNode) GetType() string { return t.Type }

// method
func NewTextNode(content string) *TextNode {
	jsonBytes, _ := json.Marshal(content)
	jsonString := string(jsonBytes)
	finalContent := strings.ReplaceAll(jsonString, "\\r\\n", "\\n")
	return &TextNode{
		Type:      TEXT_NODE,
		Content:   finalContent,
		HasEnding: false,
	}
}

func NewElementNode(tagName string) *ElementNode {
	if tagName == "" {
		panic("Empty tag was passed") //TODO: Replace with an err return
	}

	return &ElementNode{
		Type:        ELEMENT_NODE,
		Tag:         tagName,
		Data:        make(map[string]Data),
		Children:    make([]Node, 0),
		HasEnding:   false,
		IsKnownVoid: IsVoidTag(tagName),
	}
}
