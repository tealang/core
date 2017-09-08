package runtime

import "fmt"

type Evaluable interface {
	Eval(c *Context) (Value, error)
}

type Signature struct {
	Expected []Value
	Function Evaluable
}

func (sign *Signature) Match(args []Value) ([]Value, error) {
	expected, got := len(sign.Expected), len(args)
	if expected < got {
		return nil, ArgumentException(expected, got)
	}

	matched := make([]Value, expected)
	for i := range sign.Expected {
		if got > i {
			if !args[i].Type.KindOf(sign.Expected[i].Type) {
				return nil, ArgumentCastException(sign.Expected[i].Type, args[i].Type)
			}
			casted, err := sign.Expected[i].Type.Cast(args[i])
			if err != nil {
				return nil, err
			}
			matched[i] = casted
		} else if sign.Expected[i].Data != nil {
			matched[i], _ = sign.Expected[i].Type.Cast(sign.Expected[i])
		} else {
			return nil, ArgumentException(expected, got)
		}
		matched[i].Name = sign.Expected[i].Name
	}

	return matched, nil
}

func (sign *Signature) String() string {
	return fmt.Sprintf("<S {%v}>", sign.Expected)
}

type Function struct {
	Signatures []Signature
	Name       string
	Source     *Namespace
}

func (f *Function) Eval(args []Value, c *Context) (Value, error) {
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
	return Value{}, FunctionException(f.Name)
}
