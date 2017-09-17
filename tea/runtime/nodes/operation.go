package nodes

import "github.com/tealang/tea-go/tea/runtime"

// Operation calls an operator on the results of its children as arguments.
type Operation struct {
	BasicNode
	Symbol string
}

func (Operation) Name() string {
	return "Operation"
}

func (o *Operation) Eval(c *runtime.Context) (runtime.Value, error) {
	item, err := c.Namespace.Find(runtime.SearchOperator, o.Symbol)
	if err != nil {
		return runtime.Value{}, err
	}
	op, ok := item.(runtime.Operator)
	if !ok {
		return runtime.Value{}, runtime.UnexpectedItemException{Expected: runtime.Operator{}, Got: item}
	}
	args := make([]runtime.Value, len(o.Childs))
	for i, n := range o.Childs {
		v, err := n.Eval(c)
		if err != nil {
			return runtime.Value{}, err
		}
		args[i] = v
	}
	return op.Eval(c, args)
}

func NewOperation(symbol string, args ...Node) *Operation {
	return &Operation{
		BasicNode: NewBasic(args...),
		Symbol:    symbol,
	}
}
