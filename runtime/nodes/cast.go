package nodes

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/runtime"
)

// Typecast casts a value to a specific type.
type Typecast struct {
	BasicNode
	Alias string
}

func (Typecast) Name() string {
	return "Typecast"
}

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

func NewTypecast(typename string, args ...Node) *Typecast {
	return &Typecast{
		BasicNode: NewBasic(args...),
		Alias:     typename,
	}
}
