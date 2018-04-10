package nodes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/runtime"
	"github.com/tealang/core/pkg/runtime/types"
)

// Branch executes a list of conditionals until the active conditional executes successfully.
type Branch struct {
	BasicNode
}

// Name returns the name of the AST node.
func (Branch) Name() string {
	return "Branch"
}

// Eval executes the branch by iterating over all children and evaluating the conditional.
func (b *Branch) Eval(c *runtime.Context) (runtime.Value, error) {
	for _, cond := range b.Childs {
		value, err := cond.Eval(c)
		if err == nil {
			return value, nil
		} else if _, ok := err.(conditionalException); !ok {
			return runtime.Value{}, errors.Wrap(err, "failed to execute condition")
		}
	}
	return runtime.Value{}, nil
}

// NewBranch constructs a new branch node from the list of conditionals.
func NewBranch(childs ...*Conditional) *Branch {
	branch := &Branch{NewBasic()}
	for _, c := range childs {
		branch.AddBack(c)
	}
	branch.Metadata["label"] = "Branch"
	branch.Metadata["shape"] = "diamond"
	return branch
}

type conditionalException struct{}

func (c conditionalException) Error() string {
	return fmt.Sprintf("execution of conditional resulted in false")
}

type conditionalTypeException struct {
	Type *runtime.Datatype
}

func (c conditionalTypeException) Error() string {
	return fmt.Sprintf("expected value of type bool as condition result, got %s", c.Type)
}

// Conditional executes its first child, if it returns true the second child will be executed.
type Conditional struct {
	BasicNode
}

// Name returns the name of the AST node.
func (cd *Conditional) Name() string {
	return "Conditional"
}

// Eval executes the conditional by first evaluating the condition and if it results in 'true', executing the body.
func (cd *Conditional) Eval(c *runtime.Context) (runtime.Value, error) {
	condition, body := cd.Childs[0], cd.Childs[1]
	value, err := condition.Eval(c)
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "failed evaluating conditional")
	}
	if value.Type != types.Bool {
		return runtime.Value{}, conditionalTypeException{value.Type}
	}
	if !value.Data.(bool) {
		return runtime.Value{}, conditionalException{}
	}

	return c.Substitute(body.Eval)
}

// NewConditional constructs a new conditional from the given condition and body nodes.
func NewConditional(condition, body Node) *Conditional {
	cond := &Conditional{
		NewBasic(condition, body),
	}
	cond.Metadata["label"] = "Conditional"
	cond.Metadata["shape"] = "parallelogram"
	return cond
}
