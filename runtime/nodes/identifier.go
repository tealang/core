package nodes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tealang/core/runtime"
)

// Identifier is node storing a value alias that can be evaluted to retrieve the associated value.
type Identifier struct {
	BasicNode
	Alias string
}

// Name returns the name of the AST node.
func (Identifier) Name() string {
	return "Identifier"
}

// Eval retrieves the value associated with the alias in the given context namespace.
func (i *Identifier) Eval(c *runtime.Context) (runtime.Value, error) {
	item, err := c.Namespace.Find(runtime.SearchIdentifier, i.Alias)
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "failed evaluating identifier")
	}
	switch v := item.(type) {
	case runtime.Value:
		return v, nil
	default:
		return runtime.Value{}, errors.Errorf("type %T not supported", item)
	}
}

// NewIdentifier constructs a new identifier node with the given value alias.
func NewIdentifier(alias string) *Identifier {
	ident := &Identifier{
		BasicNode: NewBasic(),
		Alias:     alias,
	}
	ident.Metadata["label"] = fmt.Sprintf("Identifier (alias=%s)", alias)
	return ident
}
