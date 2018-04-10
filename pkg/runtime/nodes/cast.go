package nodes

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/runtime"
)

// Typecast casts a value to a specific type.
type Typecast struct {
	BasicNode
	Alias string
}

// Name returns the name of AST node.
func (Typecast) Name() string {
	return "Typecast"
}

// Eval executes the typecast by first looking up the target type, evaluating the children for results and casting the last result to the expected type.
func (t *Typecast) Eval(c *runtime.Context) (runtime.Value, error) {
	item, err := c.Namespace.Find(runtime.SearchDatatype, t.Alias)
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "can not typecast")
	}
	datatype, ok := item.(*runtime.Datatype)
	if !ok {
		return runtime.Value{}, errors.Errorf("expected datatype, got %s", item)
	}
	value := runtime.Value{}
	for _, n := range t.Childs {
		value, err = n.Eval(c)
		if err != nil {
			return runtime.Value{}, errors.Wrap(err, "can not typecast")
		}
	}
	return datatype.Cast(value)
}

// NewTypecast constructs a new typecast from the given typename and a list of value nodes.
func NewTypecast(typename string, args ...Node) *Typecast {
	return &Typecast{
		BasicNode: NewBasic(args...),
		Alias:     typename,
	}
}
