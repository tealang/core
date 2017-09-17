package nodes

import "github.com/tealang/tea-go/tea/runtime"

// Controller executes a range of statements and modifies the current context behavior before returning..
type Controller struct {
	BasicNode
	Behavior runtime.ContextBehavior
}

func (Controller) Name() string {
	return "Controller"
}

func (ctrl *Controller) Eval(c *runtime.Context) (runtime.Value, error) {
	var value runtime.Value
	for _, n := range ctrl.Childs {
		value, err := n.Eval(c)
		if err != nil {
			return value, err
		}
	}
	c.Behavior = ctrl.Behavior
	return value, nil
}

func NewController(behavior runtime.ContextBehavior) *Controller {
	return &Controller{
		BasicNode: NewBasic(),
		Behavior:  behavior,
	}
}
