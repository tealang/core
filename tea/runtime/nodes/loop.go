package nodes

import "github.com/tealang/tea-go/tea/runtime"

// Loop executes the conditional over and over, as long as the condition is true.
type Loop struct {
	Conditional
}

func (Loop) Name() string {
	return "Loop"
}

func (l *Loop) Eval(c *runtime.Context) (runtime.Value, error) {
	value, err := l.Conditional.Eval(c)
	_, ok := err.(ConditionalException)
	for !ok {
		switch c.Behavior {
		case runtime.BehaviorReturn:
			return value, nil
		case runtime.BehaviorBreak:
			c.Behavior = runtime.BehaviorDefault
			return runtime.Value{}, nil
		default:
			value, err = l.Conditional.Eval(c)
		}
	}
	return value, nil
}

func NewLoop(condition, body Node) *Loop {
	cond := NewConditional(condition, body)
	return &Loop{*cond}
}
