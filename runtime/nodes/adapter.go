package nodes

import "github.com/tealang/core/runtime"

// Adapter allows a Go function to be used as a tree node.
type Adapter struct {
	BasicNode
	Func func(c *runtime.Context) (runtime.Value, error)
}

func (Adapter) Name() string {
	return "Adapter"
}

func (a *Adapter) Eval(c *runtime.Context) (runtime.Value, error) {
	return a.Func(c)
}

func NewAdapter(f func(c *runtime.Context) (runtime.Value, error)) *Adapter {
	return &Adapter{
		BasicNode: NewBasic(),
		Func:      f,
	}
}
