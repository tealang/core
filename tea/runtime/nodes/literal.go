package nodes

import "github.com/tealang/tea-go/tea/runtime"

// Literal returns a constant value when evaluated.
type Literal struct {
	BasicNode
	Value runtime.Value
}

func (Literal) Name() string {
	return "Literal"
}

func (l *Literal) Eval(c *runtime.Context) (runtime.Value, error) {
	return l.Value, nil
}
