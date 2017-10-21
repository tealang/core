package functions

import (
	"fmt"

	"github.com/tealang/core/runtime"
	"github.com/tealang/core/runtime/nodes"
	"github.com/tealang/core/runtime/types"
)

func loadTypeof(c *runtime.Context) {
	typeOfSignature := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name: "data",
				Type: types.Any,
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			item, _ := c.Namespace.Find(runtime.SearchIdentifier, "data")
			value := item.(runtime.Value)
			if value.Typeflag != nil {
				return runtime.Value{
					Type: types.String,
					Data: value.Typeflag.Name,
				}, nil
			}
			return runtime.Value{
				Type: types.String,
				Data: value.Type.Name,
			}, nil
		}),
		Returns: runtime.Value{
			Type: types.String,
		},
	}
	typeOfFunction := runtime.Function{
		Signatures: []runtime.Signature{typeOfSignature},
		Source:     nil,
	}
	typeof := runtime.Value{
		Type:     types.Function,
		Data:     typeOfFunction,
		Name:     "typeof",
		Constant: true,
	}
	c.Namespace.Store(typeof)
}

func loadPrint(c *runtime.Context) {
	printSignature := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name: "text",
				Type: types.Any,
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			value, _ := c.Namespace.Find(runtime.SearchIdentifier, "text")
			fmt.Println(value.(runtime.Value).Data)
			return runtime.Value{}, nil
		}),
	}
	printFunction := runtime.Function{
		Signatures: []runtime.Signature{printSignature},
		Source:     nil,
	}
	print := runtime.Value{
		Name: "print",
		Type: types.Function,
		Data: printFunction,
	}
	c.Namespace.Store(print)
}

func Load(c *runtime.Context) {
	loadPrint(c)
	loadTypeof(c)
}
