package html

import (
	"errors"
	"strings"
	"sync"
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

func (t *TextNode) GetType() string { return t.Type }

type Data struct {
	Value     string
	IsPresent bool
}
type ElementNode struct {
	Type        string
	Tag         string
	Children    []Node
	IsKnownVoid bool
	HasEnding   bool

	data   map[string]Data
	rwLock sync.RWMutex
}

func (t *ElementNode) GetType() string { return t.Type }

func (eN *ElementNode) GetData(key string) (Data, bool) {
	eN.rwLock.RLock()
	value, exists := eN.data[key]
	eN.rwLock.RUnlock()

	return value, exists
}

func (eN *ElementNode) SetData(key string, value Data) {
	eN.rwLock.Lock()
	defer eN.rwLock.Unlock()
	eN.data[key] = value
}

// method
func NewTextNode(content string) *TextNode {
	finalContent := strings.ReplaceAll(content, "\r\n", "\n")
	return &TextNode{
		Type:      TEXT_NODE,
		Content:   finalContent,
		HasEnding: false,
	}
}

func NewElementNode(tagName string) (*ElementNode, error) {
	if tagName == "" {
		return nil, errors.New("empty tag was passed")
	}

	return &ElementNode{
		Type:        ELEMENT_NODE,
		Tag:         tagName,
		Children:    make([]Node, 0),
		HasEnding:   false,
		IsKnownVoid: IsVoidTag(tagName),
		data:        make(map[string]Data),
		rwLock:      sync.RWMutex{},
	}, nil
}
