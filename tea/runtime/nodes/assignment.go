package nodes

import "github.com/tealang/tea-go/tea/runtime"

// Assignment assigns one or multiple existing variables a new value (for each variable).
type Assignment struct {
	BasicNode
	Alias []string
}

func (Assignment) Name() string {
	return "Assignment"
}

func (a *Assignment) Eval(c *runtime.Context) (runtime.Value, error) {
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
		if err = c.Namespace.Update(value.Rename(a.Alias[i])); err != nil {
			return runtime.Value{}, err
		}
	}
	return value, nil
}

func NewMultiAssignment(alias []string, values ...Node) *Assignment {
	return &Assignment{
		BasicNode: NewBasic(values...),
		Alias:     alias,
	}
}

func NewAssignment(alias string, value Node) *Assignment {
	return &Assignment{
		BasicNode: NewBasic(value),
		Alias:     []string{alias},
	}
}
