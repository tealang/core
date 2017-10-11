package functions

import (
	"fmt"

	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

var (
	printSignature runtime.Signature
	printFunction  runtime.Function
	Print          runtime.Value
	meowSignature  runtime.Signature
	meowFunction   runtime.Function
	Meow           runtime.Value
)

func init() {
	printSignature = runtime.Signature{
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
	meowSignature = runtime.Signature{
		Expected: []runtime.Value{
			{
				Name: "separator",
				Type: types.String,
			},
			{
				Name: "values",
				Type: types.Any,
			},
			{
				Name: "x",
				Type: types.Any,
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			return runtime.Value{
				Type: types.String,
				Data: "meow",
			}, nil
		}),
	}
	meowFunction = runtime.Function{
		Signatures: []runtime.Signature{meowSignature},
		Source:     nil,
	}
	Meow = runtime.Value{
		Name: "meow",
		Type: types.Function,
		Data: meowFunction,
	}
	printFunction = runtime.Function{
		Signatures: []runtime.Signature{printSignature},
		Source:     nil,
	}
	Print = runtime.Value{
		Name: "print",
		Type: types.Function,
		Data: printFunction,
	}
}

func Load(c *runtime.Context) {
	c.Namespace.Store(Print)
	c.Namespace.Store(Meow)
}
