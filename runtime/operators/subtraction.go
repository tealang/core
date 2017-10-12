package operators

import (
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

func loadSubtraction(c *runtime.Context) {
	minusFloat := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Float,
				Constant: true,
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			var (
				identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
				a         = identA.(runtime.Value)
			)
			return runtime.Value{
				Type: types.Float,
				Data: -a.Data.(float64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Float},
	}
	minusInteger := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Integer,
				Constant: true,
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			var (
				identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
				a         = identA.(runtime.Value)
			)
			return runtime.Value{
				Type: types.Integer,
				Data: -a.Data.(int64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Integer},
	}
	subFloat := runtime.Signature{
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
			return runtime.Value{
				Type: types.Float,
				Data: a.Data.(float64) - b.Data.(float64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Float},
	}
	subInteger := runtime.Signature{
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
				Data: a.Data.(int64) - b.Data.(int64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Integer},
	}
	subFunction := runtime.Function{
		Signatures: []runtime.Signature{
			subInteger,
			subFloat,
			minusFloat,
			minusInteger,
		},
		Source: nil,
	}
	sub := runtime.Operator{
		Function: subFunction,
		Symbol:   "-",
		Constant: true,
	}
	c.Namespace.Store(sub)
}
