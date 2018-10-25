package nodes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/runtime"
)

// Assignment assigns one or multiple existing variables a new value (for each variable).
type Assignment struct {
	BasicNode
	Alias  []string
	Operator string
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
		result runtime.Value
		err   error
	)
	// Step 1: generate values
	results := make([]runtime.Value, len(a.Childs))
	for i, node := range a.Childs {
		result, err = node.Eval(c)
		if err != nil {
			return runtime.Value{}, errors.Wrap(err, "failed to assign values")
		}
		results[i] = result
	}
	// Step 2: lookup operation if necessary
	var operation *runtime.Operator
	if a.Operator != "" {
		item, err := c.Namespace.Find(runtime.SearchOperator, a.Operator)
		if err != nil {
			return runtime.Value{}, errors.Wrap(err, "missing assignment operator")
		}
		op, ok := item.(runtime.Operator)
		if !ok {
			return runtime.Value{}, errors.Errorf("expected operator, got %v", op)
		}
		operation = &op
	}

	// Step 3: store them
	for i, value := range results {
		result = value
		if operation != nil {
			item, err := c.Namespace.Find(runtime.SearchIdentifier, a.Alias[i])
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "failed to assign value")
			}
			result, err = operation.Eval(c, []runtime.Value{item.(runtime.Value), value})
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "failed to assign value")
			}
		}
		if err := c.Namespace.Update(result.Rename(a.Alias[i])); err != nil {
			return runtime.Value{}, errors.Wrap(err, "failed to assign value")
		}
	}
	return result, nil
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
