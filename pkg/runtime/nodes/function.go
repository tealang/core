package nodes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/runtime"
	"github.com/tealang/core/pkg/runtime/types"
)

// FunctionCall calls a function with the evaluated children as parameters.
type FunctionCall struct {
	BasicNode
	Alias string
}

// Name returns the name of the AST node.
func (FunctionCall) Name() string {
	return "FunctionCall"
}

// Eval calls the target function by first evaluating all its children and then feeding them as parameters to the function.
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

// NewFunctionCall constructs a new function call of the given function alias.
func NewFunctionCall(alias string, args ...Node) *FunctionCall {
	call := &FunctionCall{
		BasicNode: NewBasic(args...),
		Alias:     alias,
	}
	call.Metadata["label"] = fmt.Sprintf("Call (name='%s')", alias)
	return call
}

// FunctionLiteral generates on evaluation a new function with parameters, return type and function body.
type FunctionLiteral struct {
	BasicNode
	Args    []*Type
	Returns *Type
}

func (literal *FunctionLiteral) buildSignature(c *runtime.Context) (runtime.Signature, error) {
	// load arg types
	args := make([]runtime.Value, len(literal.Args))
	for i, arg := range literal.Args {
		value, err := arg.Eval(c)
		if err != nil {
			return runtime.Signature{}, errors.Wrap(err, "could not build signature")
		}
		args[i] = value
	}
	returns := runtime.Value{}
	// load return types
	if literal.Returns != nil {
		value, err := literal.Returns.Eval(c)
		if err != nil {
			return runtime.Signature{}, errors.Wrap(err, "failed generating return type")
		}
		returns = value
	}
	signature := runtime.NewSignature(returns, literal.Childs[0], args)
	return signature, nil
}

// Eval executes the literal and generates a single-signature function value.
func (literal *FunctionLiteral) Eval(c *runtime.Context) (runtime.Value, error) {
	signature, err := literal.buildSignature(c)
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "failed evaluating function literal")
	}
	function := runtime.NewFunction(c.Namespace, signature)
	return runtime.Value{
		Typeflag: runtime.T(types.Function),
		Data:     function,
	}, nil
}

// Name returns the name of the AST node.
func (FunctionLiteral) Name() string {
	return "FunctionLiteral"
}

// NewFunctionLiteral constructs a new function literal from the given body, return type and parameters.
func NewFunctionLiteral(body Node, returns *Type, args ...*Type) *FunctionLiteral {
	lit := &FunctionLiteral{
		BasicNode: NewBasic(body),
		Returns:   returns,
		Args:      args,
	}
	typenames := make([]string, len(args))
	for i, a := range args {
		typenames[i] = a.Tree.String()
	}
	if returns != nil {
		lit.Metadata["label"] = fmt.Sprintf("Function %s -> %s", typenames, returns.Tree)
	} else {
		lit.Metadata["label"] = fmt.Sprintf("Function %s", typenames)
	}
	return lit
}

// OperatorDefinition is an extended function literal that also has a symbol associated to it.
type OperatorDefinition struct {
	FunctionLiteral
	Symbol string
}

// Name returns the name of the AST node.
func (OperatorDefinition) Name() string {
	return "OperatorDefinition"
}

// Eval generates a function literal and stores in the context namespace as an operator.
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
		Typeflag: runtime.T(types.Function),
		Data:     function,
	}, nil
}

// NewOperatorDefinition constructs a new operator defitinition with the associated symbol.
func NewOperatorDefinition(symbol string, body Node, returns *Type, args ...*Type) *OperatorDefinition {
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
		types[i] = a.Tree.String()
	}
	if returns != nil {
		def.Metadata["label"] = fmt.Sprintf("Define %s as %s -> %s", symbol, types, returns.Tree)
	} else {
		def.Metadata["label"] = fmt.Sprintf("Define %s as %s -> ()", symbol, types)
	}
	return def
}
