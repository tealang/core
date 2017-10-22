package runtime

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
		return nil, errors.Errorf("unknown signature, expected %d args, got %d", expected, got)
	}

	matched := make([]Value, expected)
	for i := range sign.Expected {
		if got > i {
			if !args[i].Type.KindOf(sign.Expected[i].Type) {
				return nil, errors.Errorf("unknown signature, expected type %s for argument %d, got %s", sign.Expected[i].Type, i, args[i].Type)
			}
			casted, err := sign.Expected[i].Type.Cast(args[i])
			if err != nil {
				return nil, err
			}
			matched[i] = casted
		} else if sign.Expected[i].Data != nil {
			matched[i], _ = sign.Expected[i].Type.Cast(sign.Expected[i])
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
	if sign.Returns.Type != nil {
		return fmt.Sprintf("(%s) -> %s", strings.Join(items, ","), sign.Returns.Type.Name)
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
			if sign.Returns.Type != value.Type {
				return Value{}, errors.Errorf("expected return type %s, got %s", sign.Returns.Type, value.Type)
			}
			return value, nil
		})
	}
	return Value{}, errors.New("no matching signature found")
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
