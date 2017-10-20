package nodes

import (
	"fmt"

	"github.com/tealang/core/runtime"
)

type UnsupportedTypeException struct {
	Item interface{}
}

func (e UnsupportedTypeException) Error() string {
	return fmt.Sprintf("UnsupportedTypeException: Type %T not supported", e.Item)
}

type Identifier struct {
	BasicNode
	Alias string
}

func (Identifier) Name() string {
	return "Identifier"
}

func (i *Identifier) Eval(c *runtime.Context) (runtime.Value, error) {
	item, err := c.Namespace.Find(runtime.SearchIdentifier, i.Alias)
	if err != nil {
		return runtime.Value{}, err
	}
	switch v := item.(type) {
	case runtime.Value:
		return v, nil
	default:
		return runtime.Value{}, UnsupportedTypeException{v}
	}
}

func NewIdentifier(alias string) *Identifier {
	ident := &Identifier{
		BasicNode: NewBasic(),
		Alias:     alias,
	}
	ident.Metadata["label"] = fmt.Sprintf("Identifier (alias=%s)", alias)
	return ident
}
