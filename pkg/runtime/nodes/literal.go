package nodes

import (
	"fmt"

	"github.com/tealang/core/pkg/runtime"
)

// Literal returns a constant value when evaluated.
type Literal struct {
	BasicNode
	Value runtime.Value
}

// Name returns the name of the AST node.
func (Literal) Name() string {
	return "Literal"
}

// Eval returns the value of the literal.
func (l *Literal) Eval(c *runtime.Context) (runtime.Value, error) {
	return l.Value, nil
}

// NewLiteral constructs a literal with the given value.
func NewLiteral(value runtime.Value) *Literal {
	lit := &Literal{
		BasicNode: NewBasic(),
		Value:     value,
	}
	lit.Metadata["label"] = fmt.Sprintf("Literal (value='%s')", value)
	return lit
}
