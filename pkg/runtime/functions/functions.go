package functions

import (
	"fmt"
	"os"

	"github.com/tealang/core/pkg/runtime"
	"github.com/tealang/core/pkg/runtime/nodes"
	"github.com/tealang/core/pkg/runtime/types"
)

func loadStructof(c *runtime.Context) {
	structOfSignature := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "data",
				Typeflag: runtime.T(types.Any),
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			item, _ := c.Namespace.Find(runtime.SearchIdentifier, "data")
			value := item.(runtime.Value)
			return runtime.Value{
				Typeflag: runtime.T(types.String),
				Data: fmt.Sprintf("%#v", value),
			}, nil
		}),
		Returns: runtime.Value{
			Typeflag: runtime.T(types.String),
		},
	}
	structOfFunction := runtime.Function{
		Signatures: []runtime.Signature{structOfSignature},
		Source: nil,
	}
	structof := runtime.Value{
		Typeflag: runtime.T(types.Function),
		Data: structOfFunction,
		Name: "structof",
		Constant: true,
	}
	c.Namespace.Store(structof)
}

func loadTypeof(c *runtime.Context) {
	typeOfSignature := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "data",
				Typeflag: runtime.T(types.Any),
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			item, _ := c.Namespace.Find(runtime.SearchIdentifier, "data")
			value := item.(runtime.Value)
			return runtime.Value{
				Typeflag: runtime.T(types.String),
				Data:     value.Typeflag.String(),
			}, nil
		}),
		Returns: runtime.Value{
			Typeflag: runtime.T(types.String),
		},
	}
	typeOfFunction := runtime.Function{
		Signatures: []runtime.Signature{typeOfSignature},
		Source:     nil,
	}
	typeof := runtime.Value{
		Typeflag: runtime.T(types.Function),
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
				Name:     "text",
				Typeflag: runtime.T(types.Any),
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			value, _ := c.Namespace.Find(runtime.SearchIdentifier, "text")
			fmt.Fprintln(os.Stdout, value.(runtime.Value).Data)
			return runtime.Value{}, nil
		}),
	}
	printFunction := runtime.Function{
		Signatures: []runtime.Signature{printSignature},
		Source:     nil,
	}
	print := runtime.Value{
		Name:     "print",
		Typeflag: runtime.T(types.Function),
		Data:     printFunction,
	}
	c.Namespace.Store(print)
}

// Load loads language-level function into the context namespace.
func Load(c *runtime.Context) {
	loadPrint(c)
	loadTypeof(c)
	loadStructof(c)
}
