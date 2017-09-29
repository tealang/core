package nodes

import "github.com/tealang/tea-go/runtime"

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
		return runtime.Value{}, err
	}
	datatype, ok := item.(*runtime.Datatype)
	if !ok {
		return runtime.Value{}, runtime.UnexpectedItemException{Expected: new(runtime.Datatype), Got: item}
	}
	value := runtime.Value{}
	for _, n := range t.Childs {
		value, err = n.Eval(c)
		if err != nil {
			return runtime.Value{}, err
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
