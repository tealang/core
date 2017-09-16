package nodes

import "github.com/tealang/tea-go/tea/runtime"

// Declaration stores one or multiple initialized variables in the active namespace.
type Declaration struct {
	BasicNode
	Alias []string
}

func (Declaration) Name() string {
	return "Declaration"
}

func (a *Declaration) Eval(c *runtime.Context) (runtime.Value, error) {
	if len(a.Childs) != len(a.Alias) {
		return runtime.Value{}, runtime.AssignmentMismatchException{}
	}
	var (
		value runtime.Value
		err   error
	)
	for i, n := range a.Childs {
		value, err = n.Eval(c)
		if err != nil {
			return runtime.Value{}, err
		}
		if err = c.Namespace.Store(value.Rename(a.Alias[i])); err != nil {
			return runtime.Value{}, err
		}
	}
	return value, nil
}

func NewDeclaration(alias []string, values ...Node) *Declaration {
	return &Declaration{
		BasicNode: NewBasic(values...),
		Alias:     alias,
	}
}
