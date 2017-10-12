package operators

import (
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

func loadAddition(c *runtime.Context) {
	joinString := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.String,
				Constant: true,
			},
			{
				Name:     "b",
				Type:     types.String,
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
				Type: types.String,
				Data: a.Data.(string) + b.Data.(string),
			}, nil
		}),
		Returns: runtime.Value{Type: types.String},
	}
	plusFloat := runtime.Signature{
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
				Data: a.Data.(float64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Float},
	}
	plusInteger := runtime.Signature{
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
				Data: a.Data.(int64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Integer},
	}
	addFloat := runtime.Signature{
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
				Data: a.Data.(float64) + b.Data.(float64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Float},
	}
	addInteger := runtime.Signature{
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
	addFunction := runtime.Function{
		Signatures: []runtime.Signature{
			joinString,
			addInteger,
			addFloat,
			plusFloat,
			plusInteger,
		},
		Source: nil,
	}
	add := runtime.Operator{
		Function: addFunction,
		Symbol:   "+",
		Constant: true,
	}
	c.Namespace.Store(add)
}
