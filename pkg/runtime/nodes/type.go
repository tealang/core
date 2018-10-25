package nodes

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/runtime"
	"strings"
)

type Typetree struct {
	Name string
	Params []Typetree
}

func (tt Typetree) String() string {
	if len(tt.Params) > 0 {
		params := make([]string, len(tt.Params))
		for i := range tt.Params {
			params[i] = tt.Params[i].String()
		}
		return tt.Name + "<" + strings.Join(params, ", ") + ">"
	}
	return tt.Name
}

type Type struct {
	BasicNode
	Tree Typetree
}

func (Type) Name() string {
	return "Type"
}

func (t *Type) build(tree Typetree, c *runtime.Context) (runtime.Typeflag, error) {
	base, err := c.Namespace.Find(runtime.SearchDatatype, tree.Name)
	if err != nil {
		return runtime.Typeflag{}, errors.Errorf("type '%s' not found", tree.Name)
	}
	params := make([]runtime.Typeflag, len(tree.Params))
	for i := range tree.Params {
		params[i], err = t.build(tree.Params[i], c)
		if err != nil {
			return runtime.Typeflag{}, err
		}
	}
	return runtime.Typeflag{
		Type: base.(*runtime.Datatype),
		Params: params,
	}, nil
}

func (t *Type) Eval(c *runtime.Context) (runtime.Value, error) {
	typeflag, err := t.build(t.Tree, c)
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "can not build typeflag")
	}
	result, err := typeflag.Cast(runtime.Value{})
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "type has no nil value")
	}
	for i := range t.Childs {
		result, err = t.Childs[i].Eval(c)
		if err != nil {
			return runtime.Value{}, errors.Wrap(err, "can not eval")
		}
		result, err = typeflag.Cast(result)
		if err != nil {
			return runtime.Value{}, errors.Wrap(err, "can not cast")
		}
	}
	return result, nil
}

func NewType(tree Typetree, args ...Node) *Type {
	typed := &Type{
		BasicNode: NewBasic(args...),
		Tree: tree,
	}
	typed.Metadata["label"] = fmt.Sprintf("Type (%s)", tree)
	return typed
}