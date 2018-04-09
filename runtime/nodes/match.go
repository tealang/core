package nodes

import (
	"github.com/tealang/core/runtime"
)

type Match struct {
	BasicNode
}

func (Match) Name() string {
	return "Match"
}

func (m *Match) Graphviz(uid string) []string {
	m.Metadata["label"] = "Match"
	return m.BasicNode.Graphviz(uid)
}

func (m *Match) Eval(c *runtime.Context) (runtime.Value, error) {
	return c.Substitute(func(c *runtime.Context) (runtime.Value, error) {
		var result runtime.Value
		match, err := m.Childs[0].Eval(c)
		if err != nil {
			return runtime.Value{}, err
		}
		defer func() {
			if c.Behavior == runtime.BehaviorFallthrough {
				c.Behavior = runtime.BehaviorDefault
			}
		}()
		for _, n := range m.Childs[1:] {
			switch c.Behavior {
			case runtime.BehaviorFallthrough:
				result, err = n.Eval(c)
			default:
				switch n := n.(type) {
				case *Case:
					result, err = n.EvalCompare(match, c)
				default:
					result, err = n.Eval(c)
				}
			}
			if err != nil {
				if _, ok := err.(conditionalException); !ok {
					return runtime.Value{}, err
				}
				continue
			}
			if c.Behavior != runtime.BehaviorFallthrough {
				return result, err
			}
		}
		return result, nil
	})
}

// NewMatch instantiates a new 'match' node.
func NewMatch(match Node, to ...Node) *Match {
	to = append([]Node{match}, to...)
	return &Match{
		BasicNode: NewBasic(to...),
	}
}

type Case struct {
	BasicNode
}

func (Case) Name() string {
	return "Case"
}

func (c *Case) Graphviz(uid string) []string {
	c.Metadata["label"] = "Case"
	return c.BasicNode.Graphviz(uid)
}

func (c *Case) EvalCompare(match runtime.Value, ctx *runtime.Context) (runtime.Value, error) {
	value, err := c.Childs[0].Eval(ctx)
	if err != nil {
		return runtime.Value{}, err
	}
	if !value.EqualTo(match) {
		return runtime.Value{}, conditionalException{}
	}
	return c.Childs[1].Eval(ctx)
}

func (c *Case) Eval(ctx *runtime.Context) (runtime.Value, error) {
	return c.Childs[1].Eval(ctx)
}

func NewCase(value Node, body Node) *Case {
	return &Case{
		BasicNode: NewBasic(value, body),
	}
}
