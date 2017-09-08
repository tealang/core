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

type Sequence struct {
	BasicNode
	Substitute bool
}

func (n *Sequence) Name() string {
	return "Sequence"
}

func (n *Sequence) Eval(c *runtime.Context) (runtime.Value, error) {
	var parent *runtime.Namespace
	if n.Substitute {
		c.Namespace, parent = runtime.NewNamespace(c.Namespace), c.Namespace
		defer func() { c.Namespace = parent }()
	}

	for _, node := range n.Childs {
		c.Behavior = runtime.BehaviorDefault
		value, err := node.Eval(c)
		if err != nil {
			return value, err
		}
		if c.Behavior != runtime.BehaviorDefault {
			return value, nil
		}
	}

	return runtime.Value{}, nil
}

func NewSequenceNode(substitute bool, childs ...Node) *Sequence {
	return &Sequence{
		BasicNode:  NewBasicNode(childs...),
		Substitute: substitute,
	}
}
