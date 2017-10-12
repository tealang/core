package operators

import (
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

func loadDivision(c *runtime.Context) {
	divFloat := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Float,
				Constant: true,
			},
			{
				Name:     "b",
				Type:     types.Float,
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
			if b.Data.(float64) == 0.0 {
				return runtime.Value{}, runtime.RuntimeException{Message: "Cannot divide by 0"}
			}
			return runtime.Value{
				Type: types.Float,
				Data: a.Data.(float64) / b.Data.(float64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Float},
	}
	divInteger := runtime.Signature{
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
			if b.Data.(int64) == 0 {
				return runtime.Value{}, runtime.RuntimeException{Message: "Cannot divide by 0"}
			}
			return runtime.Value{
				Type: types.Integer,
				Data: a.Data.(int64) / b.Data.(int64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Integer},
	}
	divFunction := runtime.Function{
		Signatures: []runtime.Signature{
			divInteger,
			divFloat,
		},
		Source: nil,
	}
	div := runtime.Operator{
		Function: divFunction,
		Symbol:   "/",
		Constant: true,
	}
	c.Namespace.Store(div)
}
