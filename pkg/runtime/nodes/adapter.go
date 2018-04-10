package nodes

import "github.com/tealang/core/pkg/runtime"

// Adapter allows a Go function to be used as a tree node.
type Adapter struct {
	BasicNode
	Func func(c *runtime.Context) (runtime.Value, error)
}

// Name returns the name of the AST node.
func (Adapter) Name() string {
	return "Adapter"
}

// Eval executes the encapsulated function in the call context.
func (a *Adapter) Eval(c *runtime.Context) (runtime.Value, error) {
	return a.Func(c)
}

// NewAdapter constructs a new adapter using the given function.
func NewAdapter(f func(c *runtime.Context) (runtime.Value, error)) *Adapter {
	return &Adapter{
		BasicNode: NewBasic(),
		Func:      f,
	}
}
