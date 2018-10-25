package runtime

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Evaluable can be executed and results in a value.
type Evaluable interface {
	Eval(c *Context) (Value, error)
}

// Signature stores a function signature consisting of arguments, return value and function body.
type Signature struct {
	Expected []Value
	Function Evaluable
	Returns  Value
}

// Match checks if the arguments match to the signature.
func (sign Signature) Match(args []Value) ([]Value, error) {
	expected, got := len(sign.Expected), len(args)
	if expected < got {
		return nil, errors.Errorf("too many args, expected %d args, got %d", expected, got)
	}

	matched := make([]Value, expected)
	for i := range sign.Expected {
		if got > i {
			if !args[i].Type.KindOf(sign.Expected[i].Type) {
				return nil, errors.Errorf("unknown signature, expected type %s for argument %d, got %s", sign.Expected[i].Type, i, args[i].Type)
			}
			casted, err := sign.Expected[i].Cast(args[i])
			if err != nil {
				return nil, errors.Wrap(err, "signature not matching")
			}
			matched[i] = casted
		} else if sign.Expected[i].Data != nil {
			matched[i], _ = sign.Expected[i].Cast(sign.Expected[i])
		} else {
			return nil, errors.Errorf("missing args, expected %d, got %d", expected, got)
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
		return fmt.Sprintf("(%s) -> %s", strings.Join(items, ","), sign.Returns.Typeflag)
	}
	return fmt.Sprintf("(%s)", strings.Join(items, ","))
}

// NewSignature creates a new signature.
func NewSignature(returns Value, executes Evaluable, args []Value) Signature {
	return Signature{
		Returns:  returns,
		Function: executes,
		Expected: args,
	}
}

// Function is a collection of signatures with a common source namespace.
type Function struct {
	Signatures []Signature
	Source     *Namespace
}

// Eval executes the function, searching and executing a matching signature.
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
				return Value{}, errors.Wrap(err, "failed to evaluate")
			}
			if sign.Returns.Type != nil && !value.Type.KindOf(sign.Returns.Type) {
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

// NewFunction creates a new function using the given source namespace and signatures.
func NewFunction(source *Namespace, signatures ...Signature) Function {
	return Function{
		Source:     source,
		Signatures: signatures,
	}
}
