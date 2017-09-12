package nodes

import "github.com/tealang/tea-go/tea/runtime"

// Sequence executes a list of nodes as long as the behavior is default.
type Sequence struct {
	BasicNode
	Substitute bool
}

func (Sequence) Name() string {
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

func NewSequence(substitute bool, childs ...Node) *Sequence {
	return &Sequence{
		BasicNode:  NewBasic(childs...),
		Substitute: substitute,
	}
}
