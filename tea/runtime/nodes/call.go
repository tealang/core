package nodes

import (
	"github.com/tealang/tea-go/tea/runtime"
	"github.com/tealang/tea-go/tea/stdlib/types"
)

type FunctionCall struct {
	BasicNode
	Alias string
}

func (FunctionCall) Name() string {
	return "FunctionCall"
}

func (call *FunctionCall) Eval(c *runtime.Context) (runtime.Value, error) {
	item, err := c.Namespace.Find(runtime.SearchIdentifier, call.Alias)
	if err != nil {
		return runtime.Value{}, err
	}
	value, ok := item.(runtime.Value)
	if err != nil {
		return runtime.Value{}, runtime.UnexpectedItemException{Expected: runtime.Value{}, Got: item}
	}
	if !value.Type.KindOf(types.Function) {
		return runtime.Value{}, runtime.UncallableTypeException{Type: value.Type}
	}
	callable, ok := value.Data.(runtime.Function)
	if !ok {
		return runtime.Value{}, runtime.UncallableTypeException{Type: value.Type}
	}
	values := make([]runtime.Value, len(call.Childs))
	for i, n := range call.Childs {
		v, err := n.Eval(c)
		if err != nil {
			return runtime.Value{}, err
		}
		values[i] = v
	}
	return callable.Eval(values, c)
}

func NewFunctionCall(alias string, args ...Node) *FunctionCall {
	return &FunctionCall{
		BasicNode: NewBasic(args...),
		Alias:     alias,
	}
}
