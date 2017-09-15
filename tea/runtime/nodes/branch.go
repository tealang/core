package nodes

import (
	"fmt"

	"github.com/tealang/tea-go/tea/runtime"
	"github.com/tealang/tea-go/tea/stdlib/types"
)

// Branch executes a list of conditionals until the active conditional executes successfully.
type Branch struct {
	BasicNode
}

func (Branch) Name() string {
	return "Branch"
}

func (b *Branch) Eval(c *runtime.Context) (runtime.Value, error) {
	for _, cond := range b.Childs {
		value, err := cond.Eval(c)
		if err == nil {
			return value, nil
		} else if _, ok := err.(ConditionalException); !ok {
			return value, err
		}
	}
	return runtime.Value{}, nil
}

func NewBranch(childs ...*Conditional) *Branch {
	branch := &Branch{NewBasic()}
	for _, c := range childs {
		branch.AddBack(c)
	}
	return branch
}

type ConditionalException struct{}

func (c ConditionalException) Error() string {
	return fmt.Sprintf("ConditionalException: Execution resulted in false")
}

type ConditionalTypeException struct {
	Type *runtime.Datatype
}

func (c ConditionalTypeException) Error() string {
	return fmt.Sprintf("ConditionalTypeException: Expected type bool, got %s", c.Type)
}

// Conditional executes its first child, if it returns true the second child will be executed.
type Conditional struct {
	BasicNode
}

func (cd *Conditional) Name() string {
	return "Conditional"
}

func (cd *Conditional) Eval(c *runtime.Context) (runtime.Value, error) {
	condition := cd.Childs[0]
	body := cd.Childs[1]

	value, err := condition.Eval(c)
	if err != nil {
		return value, err
	}
	if value.Type != types.Bool {
		return value, ConditionalTypeException{value.Type}
	}
	if !value.Data.(bool) {
		return value, ConditionalException{}
	}

	return c.Substitute(body.Eval)
}

func NewConditional(condition, body Node) *Conditional {
	return &Conditional{
		NewBasic(condition, body),
	}
}
