package nodes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tealang/core/runtime"
)

// Assignment assigns one or multiple existing variables a new value (for each variable).
type Assignment struct {
	BasicNode
	Alias []string
}

// Name returns the name of the AST node.
func (Assignment) Name() string {
	return "Assignment"
}

// Graphviz generates a graph representation of the node.
func (a *Assignment) Graphviz(uid string) []string {
	a.Metadata["label"] = fmt.Sprintf("Assignment (%s)", a.Alias)
	return a.BasicNode.Graphviz(uid)
}

// Eval executes the assignment by evaluating the value nodes and assigning the results to the values in the context namespace.
func (a *Assignment) Eval(c *runtime.Context) (runtime.Value, error) {
	if len(a.Childs) != len(a.Alias) {
		return runtime.Value{}, errors.Errorf("can not assign %d values to %d names", len(a.Childs), len(a.Alias))
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
			return runtime.Value{}, errors.Wrap(err, "failed to assign values")
		}
		results[i] = value
	}
	// Step 2: store them
	for i, value := range results {
		if err = c.Namespace.Update(value.Rename(a.Alias[i])); err != nil {
			return runtime.Value{}, errors.Wrap(err, "failed to assign values")
		}
	}
	return value, nil
}

// NewMultiAssignment constructs a new tuple assignment.
func NewMultiAssignment(alias []string, values ...Node) *Assignment {
	return &Assignment{
		BasicNode: NewBasic(values...),
		Alias:     alias,
	}
}

// NewAssignment constructs a new single-value assignment.
func NewAssignment(alias string, value Node) *Assignment {
	return NewMultiAssignment([]string{alias}, value)
}
