package main

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestLexter(t *testing.T) {
	nodes, err := Run("tests/lexer/input.html")

	if err != nil {
		t.Errorf("Failed to parse html, error: %s", err.Error())
	}

	snaps.MatchSnapshot(t, nodes)
}
