package runtime

import (
	"fmt"
	"strings"
)

type Evaluable interface {
	Eval(c *Context) (Value, error)
}

type Signature struct {
	Expected []Value
	Function Evaluable
	Returns  Value
}

func (sign Signature) Match(args []Value) ([]Value, error) {
	expected, got := len(sign.Expected), len(args)
	if expected < got {
		return nil, ArgumentException{
			Expected: expected,
			Got:      got,
		}
	}

	matched := make([]Value, expected)
	for i := range sign.Expected {
		if got > i {
			if !args[i].Type.KindOf(sign.Expected[i].Type) {
				return nil, ArgumentCastException{
					Expected: sign.Expected[i].Type,
					Got:      args[i].Type,
				}
			}
			casted, err := sign.Expected[i].Type.Cast(args[i])
			if err != nil {
				return nil, err
			}
			matched[i] = casted
		} else if sign.Expected[i].Data != nil {
			matched[i], _ = sign.Expected[i].Type.Cast(sign.Expected[i])
		} else {
			return nil, ArgumentException{
				Expected: expected,
				Got:      got,
			}
		}
		matched[i].Name = sign.Expected[i].Name
	}

	return matched, nil
}

func (sign Signature) String() string {
	items := make([]string, len(sign.Expected))
	for i, n := range sign.Expected {
		items[i] = n.VariableString()
	}
	return fmt.Sprintf("(%s)", strings.Join(items, ","))
}

func NewSignature(returns Value, executes Evaluable, args []Value) Signature {
	return Signature{
		Returns:  returns,
		Function: executes,
		Expected: args,
	}
}

type Function struct {
	Signatures []Signature
	Source     *Namespace
}

func (f Function) Eval(c *Context, args []Value) (Value, error) {
	for _, sign := range f.Signatures {
		matched, err := sign.Match(args)
		if err != nil {
			continue
		}
		return c.Substitute(func(c *Context) (Value, error) {
			c.Namespace = NewNamespace(f.Source)
			for _, arg := range matched {
				c.Namespace.Store(arg)
			}
			value, err := sign.Function.Eval(c)
			if err != nil {
				return Value{}, err
			}
			if sign.Returns.Type != nil && !value.Type.KindOf(sign.Returns.Type) {
				return Value{}, CastException{From: value.Type, To: sign.Returns.Type}
			}
			return value, nil
		})
	}
	return Value{}, FunctionException{}
}

func (f Function) String() string {
	items := make([]string, len(f.Signatures))
	for i, n := range f.Signatures {
		items[i] = n.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(items, ";"))
}

func NewFunction(source *Namespace, signatures ...Signature) Function {
	return Function{
		Source:     source,
		Signatures: signatures,
	}
}
