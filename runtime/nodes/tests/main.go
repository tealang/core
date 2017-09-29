package main

import (
	"fmt"

	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/stdlib/types"

	"github.com/tealang/tea-go/runtime"
)

func main() {
	ctx := runtime.NewContext()
	initY := nodes.NewDeclaration("y", false, nodes.NewLiteral(types.False))
	fmt.Println(initY.Eval(ctx))
	initF := nodes.NewDeclaration("f", true, nodes.NewLiteral(runtime.Value{
		Type: types.Function,
		Data: runtime.Function{
			Signatures: []runtime.Signature{runtime.Signature{
				Expected: []runtime.Value{},
				Function: nodes.NewSequence(false,
					nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
						fmt.Println("hello")
						return runtime.Value{}, nil
					}),
					nodes.NewDeclaration("y", true, nodes.NewLiteral(types.True)),
				),
			}},
			Source: ctx.Namespace,
		},
	}))
	callF := nodes.NewFunctionCall("f")
	fmt.Println(initF.Eval(ctx))
	fmt.Println(callF.Eval(ctx))
	fmt.Println(ctx.GlobalNamespace.Find(runtime.SearchIdentifier, "y"))
}
