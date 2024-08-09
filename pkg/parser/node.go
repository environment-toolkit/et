package parser

import (
	"fmt"

	"github.com/environment-toolkit/et/pkg/models"
)

// NodeType identifies the type of a parse tree node.
type NodeType int

// Type returns itself and provides an easy default implementation
// for embedding in a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeText NodeType = iota // Plain text.
	NodeVariable
)

// A Node is an element in the parse tree. The interface is trivial.
// The interface contains an unexported method so that only
// types local to this package can satisfy it.
type Node interface {
	Type() NodeType
	String() (string, error)
}

type textNode struct {
	NodeType
	Text string
}

func NewText(text string) Node {
	return &textNode{NodeText, text}
}

func (t *textNode) String() (string, error) {
	return t.Text, nil
}

type variableNode struct {
	NodeType
	Value *models.Value
}

func NewVariable(value *models.Value) Node {
	return &variableNode{NodeText, value}
}

func (t *variableNode) String() (string, error) {
	out, resolved := t.Value.Resolved()
	if !resolved {
		return "", fmt.Errorf("variable '%s' not resolved", out)
	}
	return out, nil
}
