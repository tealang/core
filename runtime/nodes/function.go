package nodes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tealang/core/runtime"
	"github.com/tealang/core/runtime/types"
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
		return runtime.Value{}, errors.Wrap(err, "undefined function")
	}
	value, ok := item.(runtime.Value)
	if err != nil {
		return runtime.Value{}, errors.Errorf("expected value, got %s", item)
	}
	if !value.Type.KindOf(types.Function) {
		return runtime.Value{}, errors.Errorf("can not call value of type %s", value.Type)
	}
	callable, ok := value.Data.(runtime.Function)
	if !ok {
		return runtime.Value{}, errors.Errorf("expected function, got %s", value.Data)
	}
	values := make([]runtime.Value, len(call.Childs))
	for i, n := range call.Childs {
		v, err := n.Eval(c)
		if err != nil {
			return runtime.Value{}, errors.Wrap(err, "can not call function")
		}
		values[i] = v
	}
	result, err := callable.Eval(c, values)
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "function call failed")
	}
	c.Behavior = runtime.BehaviorDefault
	return result, nil
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

func (literal *FunctionLiteral) buildSignature(c *runtime.Context) (runtime.Signature, error) {
	// load arg types
	args := make([]runtime.Value, len(literal.Args))
	for i, arg := range literal.Args {
		value, err := arg.Eval(c)
		if err != nil {
			return runtime.Signature{}, err
		}
		args[i] = value
	}
	returns := runtime.Value{}
	// load return types
	if literal.Returns != nil {
		value, err := literal.Returns.Eval(c)
		if err != nil {
			return runtime.Signature{}, err
		}
		returns = value
	}
	signature := runtime.NewSignature(returns, literal.Childs[0], args)
	return signature, nil
}

func (literal *FunctionLiteral) Eval(c *runtime.Context) (runtime.Value, error) {
	signature, err := literal.buildSignature(c)
	if err != nil {
		return runtime.Value{}, err
	}
	function := runtime.NewFunction(c.Namespace, signature)
	return runtime.Value{
		Type: types.Function,
		Data: function,
	}, nil
}

func (FunctionLiteral) Name() string {
	return "FunctionLiteral"
}

func NewFunctionLiteral(body Node, returns *Typecast, args ...*Typecast) *FunctionLiteral {
	lit := &FunctionLiteral{
		BasicNode: NewBasic(body),
		Returns:   returns,
		Args:      args,
	}
	types := make([]string, len(args))
	for i, a := range args {
		types[i] = a.Alias
	}
	if returns != nil {
		lit.Metadata["label"] = fmt.Sprintf("Function <%s> -> %s", types, returns.Alias)
	} else {
		lit.Metadata["label"] = fmt.Sprintf("Function <%s>", types)
	}
	return lit
}

type OperatorDefinition struct {
	FunctionLiteral
	Symbol string
}

func (OperatorDefinition) Name() string {
	return "OperatorDefinition"
}

func (definition *OperatorDefinition) Eval(c *runtime.Context) (runtime.Value, error) {
	signature, err := definition.buildSignature(c)
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "can not build signature")
	}
	function := runtime.NewFunction(c.Namespace, signature)
	operator := runtime.Operator{
		Function: function,
		Symbol:   definition.Symbol,
		Constant: true,
	}
	if err := c.Namespace.Store(operator); err != nil {
		return runtime.Value{}, errors.Wrap(err, "can not store operator")
	}
	return runtime.Value{
		Type: types.Function,
		Data: function,
	}, nil
}

func NewOperatorDefinition(symbol string, body Node, returns *Typecast, args ...*Typecast) *OperatorDefinition {
	def := &OperatorDefinition{
		FunctionLiteral: FunctionLiteral{
			BasicNode: NewBasic(body),
			Returns:   returns,
			Args:      args,
		},
		Symbol: symbol,
	}
	types := make([]string, len(args))
	for i, a := range args {
		types[i] = a.Alias
	}
	if returns != nil {
		def.Metadata["label"] = fmt.Sprintf("Define %s as <%s> -> %s", symbol, types, returns.Alias)
	} else {
		def.Metadata["label"] = fmt.Sprintf("Define %s as <%s> -> ()", symbol, types)
	}
	return def
}
