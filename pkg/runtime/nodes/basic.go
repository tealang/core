// Package nodes provides an abstract syntax tree implementation for Tealang programs.
package nodes

import (
	"fmt"

	"github.com/tealang/core/pkg/runtime"
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

// Graphviz generates a graphviz-compatible graph representation of this node and its children.
func (n *BasicNode) Graphviz(uid string) []string {
	lines := []string{fmt.Sprintf("%s [label=\"%s\" shape=%s]", uid, n.Metadata["label"], n.Metadata["shape"])}
	for i, c := range n.Childs {
		id := fmt.Sprintf("%sn%d", uid, i)
		lines = append(lines, c.Graphviz(id)...)
		lines = append(lines, fmt.Sprintf("%s -> %s", uid, id))
	}
	return lines
}

// Has checks if the node has metadata assigned to the given key.
func (n *BasicNode) Has(tag string) (string, bool) {
	v, ok := n.Metadata[tag]
	return v, ok
}

// Tag assigns metadata to a key.
func (n *BasicNode) Tag(tag, value string) {
	n.Metadata[tag] = value
}

// AddBack adds a new child at the end of the children list.
func (n *BasicNode) AddBack(child Node) {
	n.Childs = append(n.Childs, child)
}

// AddFront adds a new child at the front of the children list.
func (n *BasicNode) AddFront(child Node) {
	n.Childs = append([]Node{child}, n.Childs...)
}

// NewBasic constructs a new basic node that can not be evaluated.
func NewBasic(childs ...Node) BasicNode {
	return BasicNode{
		Childs: childs,
		Metadata: map[string]string{
			"label": "Node",
			"shape": "oval",
		},
	}
}
