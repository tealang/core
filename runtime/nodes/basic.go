// Package nodes provides an abstract syntax tree implementation for Tealang programs.
package nodes

import (
	"fmt"

	"github.com/tealang/tea-go/runtime"
)

// Node is a tree node within an abstract syntax tree.
type Node interface {
	runtime.Evaluable
	Has(tag string) (string, bool)
	Tag(tag, value string)
	Name() string
	AddFront(child Node)
	AddBack(child Node)
	Graphviz(uid string) []string
}

// BasicNode provides a basic functionality for a new node.
type BasicNode struct {
	Childs   []Node
	Metadata map[string]string
}

func (n *BasicNode) Graphviz(uid string) []string {
	lines := []string{fmt.Sprintf("%s [label=\"%s\"]", uid, n.Metadata["label"])}
	for i, c := range n.Childs {
		id := fmt.Sprintf("%sn%d", uid, i)
		lines = append(lines, c.Graphviz(id)...)
		lines = append(lines, fmt.Sprintf("%s -> %s", uid, id))
	}
	return lines
}

func (n *BasicNode) Has(tag string) (string, bool) {
	v, ok := n.Metadata[tag]
	return v, ok
}

func (n *BasicNode) Tag(tag, value string) {
	n.Metadata[tag] = value
}

func (n *BasicNode) AddBack(child Node) {
	n.Childs = append(n.Childs, child)
}

func (n *BasicNode) AddFront(child Node) {
	n.Childs = append([]Node{child}, n.Childs...)
}

func NewBasic(childs ...Node) BasicNode {
	return BasicNode{
		Childs:   childs,
		Metadata: make(map[string]string),
	}
}
