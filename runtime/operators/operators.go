package operators

import (
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

var (
	addBinary, addUnary runtime.Signature
	addFunction         runtime.Function
	Add                 runtime.Operator
)

func init() {
	addBinary = runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Integer,
				Constant: true,
			},
			{
				Name:     "b",
				Type:     types.Integer,
				Constant: true,
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			var (
				identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
				identB, _ = c.Namespace.Find(runtime.SearchIdentifier, "b")
				a         = identA.(runtime.Value)
				b         = identB.(runtime.Value)
			)
			return runtime.Value{
				Type: types.Integer,
				Data: a.Data.(int64) + b.Data.(int64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Integer},
	}
	addFunction = runtime.Function{
		Signatures: []runtime.Signature{
			addBinary,
		},
		Source: nil,
	}
	Add = runtime.Operator{
		Function: addFunction,
		Symbol:   "+",
		Constant: true,
	}
}

func Load(c *runtime.Context) {
	c.Namespace.Store(Add)
}
