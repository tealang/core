package nodes

import (
	"fmt"

	"github.com/tealang/core/runtime"
)

// Sequence executes a list of nodes as long as the behavior is default.
type Sequence struct {
	BasicNode
	Substitute bool
}

func (Sequence) Name() string {
	return "Sequence"
}

func (n *Sequence) Eval(c *runtime.Context) (runtime.Value, error) {
	var (
		parent *runtime.Namespace
		value  runtime.Value
		err    error
	)

	if n.Substitute {
		c.Namespace, parent = runtime.NewNamespace(c.Namespace), c.Namespace
		defer func() { c.Namespace = parent }()
	}
	for _, node := range n.Childs {
		c.Behavior = runtime.BehaviorDefault
		value, err = node.Eval(c)
		if err != nil {
			return value, err
		}
		if c.Behavior != runtime.BehaviorDefault {
			break
		}
	}

	return value, nil
}

func NewSequence(substitute bool, childs ...Node) *Sequence {
	seq := &Sequence{
		BasicNode:  NewBasic(childs...),
		Substitute: substitute,
	}
	seq.Metadata["label"] = fmt.Sprintf("Sequence (sub=%t)", substitute)
	seq.Metadata["shape"] = "house"
	return seq
}
