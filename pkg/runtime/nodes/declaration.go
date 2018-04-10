package nodes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/runtime"
)

// Declaration stores one or multiple initialized variables in the active namespace.
type Declaration struct {
	BasicNode
	Alias    []string
	Constant bool
}

// Graphviz generates a graphviz-compatible representation of the declaration.
func (a *Declaration) Graphviz(uid string) []string {
	a.Metadata["label"] = fmt.Sprintf("Declaration (alias=%s, constant=%t)", a.Alias, a.Constant)
	return a.BasicNode.Graphviz(uid)
}

// Name returns the name of the AST node.
func (Declaration) Name() string {
	return "Declaration"
}

// Eval executes the declaration by first retrieving the values to be assigned and then storing them in the context namespace.
func (a *Declaration) Eval(c *runtime.Context) (runtime.Value, error) {
	if len(a.Childs) != len(a.Alias) {
		return runtime.Value{}, errors.Errorf("can not declare %d values and assign to %d names", len(a.Childs), len(a.Alias))
	}
	var (
		value runtime.Value
		err   error
	)
	// Step 1: generate values
	results := make([]runtime.Value, len(a.Childs))
	for i, node := range a.Childs {
		value, err = node.Eval(c)
		if err != nil {
			return runtime.Value{}, errors.Wrap(err, "failed declaring values")
		}
		results[i] = value
	}
	// Step 2: store them
	for i, value := range results {
		if err = c.Namespace.Store(value.Rename(a.Alias[i]).Rechange(a.Constant)); err != nil {
			return runtime.Value{}, errors.Wrap(err, "failed declaring values")
		}
	}
	return value, nil
}

// NewMultiDeclaration constructs a new tuple declaration for a list of value aliases.
func NewMultiDeclaration(alias []string, constant bool, values ...Node) *Declaration {
	return &Declaration{
		BasicNode: NewBasic(values...),
		Alias:     alias,
		Constant:  constant,
	}
}

// NewDeclaration constructs a new single-value declaration.
func NewDeclaration(alias string, constant bool, value Node) *Declaration {
	return NewMultiDeclaration([]string{alias}, constant, value)
}
