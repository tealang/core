package nodes

import "github.com/tealang/tea-go/tea/runtime"

type Node interface {
	runtime.Evaluable
	Has(tag string) (string, bool)
	Tag(tag, value string)
	Name() string
	AddFront(child Node)
	AddBack(child Node)
}

type BasicNode struct {
	Childs   []Node
	Metadata map[string]string
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

func NewBasicNode(childs ...Node) BasicNode {
	return BasicNode{
		Childs:   childs,
		Metadata: make(map[string]string),
	}
}
