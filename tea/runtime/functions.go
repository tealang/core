package runtime

import "fmt"

type Evaluable interface {
	Eval(c *Context) (Value, error)
}

type Signature struct {
	Expected []Value
	Function Evaluable
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
	return fmt.Sprintf("<S {%v}>", sign.Expected)
}

func NewSignature(eval Evaluable, args ...Value) Signature {
	return Signature{
		Function: eval,
		Expected: args,
	}
}

type Function struct {
	Signatures []Signature
	Source     *Namespace
}

func (f Function) Eval(args []Value, c *Context) (Value, error) {
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
			return sign.Function.Eval(c)
		})
	}
	return Value{}, FunctionException{}
}

func NewFunction(source *Namespace, signatures ...Signature) Function {
	return Function{
		Source:     source,
		Signatures: signatures,
	}
}
