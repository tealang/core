package functions

import (
	"bufio"
	"fmt"
	"github.com/tealang/core/pkg/runtime"
	"github.com/tealang/core/pkg/runtime/nodes"
	"github.com/tealang/core/pkg/runtime/types"
	"os"
	"strings"
)

var runtimeFunctions = []func(*runtime.Context){
	loadTypeof,
	loadPrint,
	loadRead,
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

func loadRead(c *runtime.Context) {
	reader := bufio.NewReader(os.Stdin)
	readAdapter := nodes.NewAdapter(func (c *runtime.Context) (runtime.Value, error) {
		value, _ := c.Namespace.Find(runtime.SearchIdentifier, "text")
		if value != nil {
			fmt.Fprint(os.Stdout, value.(runtime.Value).Data)
		}
		input, _ := reader.ReadString('\n')
		return runtime.Value{
			Typeflag: runtime.T(types.String),
			Data: strings.TrimSuffix(input, "\n"),
		}, nil
	})
	readFunction := runtime.Function{
		Signatures: []runtime.Signature{
			{
				Expected: []runtime.Value{{Name: "text", Typeflag: runtime.T(types.String)}},
				Function: readAdapter,
			},
			{
				Expected: []runtime.Value{},
				Function: readAdapter,
			},
		},
		Source: nil,
	}
	read := runtime.Value{
		Name: "read",
		Typeflag: runtime.T(types.Function),
		Data: readFunction,
	}
	c.Namespace.Store(read)
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
	for _, f := range runtimeFunctions {
		f(c)
	}
}
