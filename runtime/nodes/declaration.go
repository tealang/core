package nodes

import (
	"fmt"

	"github.com/tealang/tea-go/runtime"
)

// Declaration stores one or multiple initialized variables in the active namespace.
type Declaration struct {
	BasicNode
	Alias    []string
	Constant bool
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
		if err = c.Namespace.Store(value.Rename(a.Alias[i]).Rechange(a.Constant)); err != nil {
			return runtime.Value{}, err
		}
	}
	return value, nil
}

func NewMultiDeclaration(alias []string, constant bool, values ...Node) *Declaration {
	decl := &Declaration{
		BasicNode: NewBasic(values...),
		Alias:     alias,
		Constant:  constant,
	}
	decl.Metadata["label"] = fmt.Sprintf("Declaration (alias=%v, constant=%t)", alias, constant)
	return decl
}

func NewDeclaration(alias string, constant bool, value Node) *Declaration {
	decl := &Declaration{
		BasicNode: NewBasic(value),
		Alias:     []string{alias},
		Constant:  constant,
	}
	decl.Metadata["label"] = fmt.Sprintf("Declaration (alias=%v, constant=%t)", alias, constant)
	return decl
}
