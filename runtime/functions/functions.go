package functions

import (
	"fmt"

	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

var (
	printString   runtime.Signature
	printFunction runtime.Function
	print         runtime.Value
	printSep      runtime.Signature
)

func init() {
	printString = runtime.Signature{
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
	printSep = runtime.Signature{
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
			fmt.Println("meow")
			return runtime.Value{}, nil
		}),
	}
	printFunction = runtime.Function{
		Signatures: []runtime.Signature{printSep, printString},
		Source:     nil,
	}
	print = runtime.Value{
		Name: "print",
		Type: types.Function,
		Data: printFunction,
	}
}

func Load(c *runtime.Context) {
	c.Namespace.Store(print)
}
