package nodes

import (
	"github.com/tealang/core/pkg/runtime"
)

// Loop executes the conditional over and over, as long as the condition is true.
type Loop struct {
	Conditional
}

// Name returns the name of the AST node.
func (Loop) Name() string {
	return "Loop"
}

// Eval executes a conditional over and over until the condition is false.
// The control flow can be manipulated using behavior control.
func (l *Loop) Eval(c *runtime.Context) (runtime.Value, error) {
	value, err := l.Conditional.Eval(c)
	_, ok := err.(conditionalException)
	for !ok {
		switch c.Behavior {
		case runtime.BehaviorReturn:
			return value, nil
		case runtime.BehaviorBreak:
			c.Behavior = runtime.BehaviorDefault
			return runtime.Value{}, nil
		default:
			value, err = l.Conditional.Eval(c)
			_, ok = err.(conditionalException)
		}
	}
	return value, nil
}

// NewLoop constructs a new loop with a condition head and body.
func NewLoop(condition, body Node) *Loop {
	cond := NewConditional(condition, body)
	loop := &Loop{*cond}
	loop.Metadata["label"] = "Loop"
	return loop
}
