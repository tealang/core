package nodes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/runtime"
)

// Sequence executes a list of nodes as long as the behavior is default.
type Sequence struct {
	BasicNode
	Substitute bool
}

// Name returns the name of the AST node.
func (Sequence) Name() string {
	return "Sequence"
}

// Eval executes the sequence by evaluating its children one by one.
// The control flow can be modified by using behavior control.
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
			return value, errors.Wrap(err, "failed evaluating sequence")
		}
		if c.Behavior != runtime.BehaviorDefault {
			break
		}
	}

	return value, nil
}

// NewSequence constructs a new sequence from the given children list.
// The sequence can also be evaluated in a substitute namespace.
func NewSequence(substitute bool, childs ...Node) *Sequence {
	seq := &Sequence{
		BasicNode:  NewBasic(childs...),
		Substitute: substitute,
	}
	seq.Metadata["label"] = fmt.Sprintf("Sequence (sub=%t)", substitute)
	seq.Metadata["shape"] = "house"
	return seq
}
