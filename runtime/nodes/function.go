package nodes

import (
	"fmt"

	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/types"
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
	return callable.Eval(c, values)
}

func NewFunctionCall(alias string, args ...Node) *FunctionCall {
	call := &FunctionCall{
		BasicNode: NewBasic(args...),
		Alias:     alias,
	}
	call.Metadata["label"] = fmt.Sprintf("Call (name='%s')", alias)
	return call
}

type FunctionLiteral struct {
	BasicNode
	Args    []*Typecast
	Returns *Typecast
}

func (literal *FunctionLiteral) Eval(c *runtime.Context) (runtime.Value, error) {
	// load arg types
	args := make([]runtime.Value, len(literal.Childs)-1)
	for i, arg := range literal.Args {
		value, err := arg.Eval(c)
		if err != nil {
			return runtime.Value{}, err
		}
		args[i] = value
	}
	returns := runtime.Value{}
	// load return types
	if literal.Returns != nil {
		value, err := literal.Returns.Eval(c)
		if err != nil {
			return runtime.Value{}, err
		}
		returns = value
	}
	signature := runtime.NewSignature(returns, literal.Childs[0], args)
	function := runtime.NewFunction(c.Namespace, signature)
	return runtime.Value{
		Type: types.Function,
		Data: function,
	}, nil
}

func NewFunctionLiteral(body, returns *Typecast, args ...*Typecast) *FunctionLiteral {
	return &FunctionLiteral{
		BasicNode: NewBasic(body),
		Returns:   returns,
		Args:      args,
	}
}
