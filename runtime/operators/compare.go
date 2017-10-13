package operators

import (
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

func loadSmaller(c *runtime.Context) {
	adapter := nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
		var (
			identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
			identB, _ = c.Namespace.Find(runtime.SearchIdentifier, "b")
			a         = identA.(runtime.Value)
			b         = identB.(runtime.Value)
		)
		b, err := a.Type.Cast(b)
		if err != nil {
			return runtime.Value{}, err
		}
		switch a.Type {
		case types.Integer:
			return runtime.Value{
				Type: types.Bool,
				Data: a.Data.(int64) < b.Data.(int64),
			}, nil
		case types.Float:
			return runtime.Value{
				Type: types.Bool,
				Data: a.Data.(float64) < b.Data.(float64),
			}, nil
		}
		return runtime.Value{}, runtime.NewNotApplicableException("<", a.Type, b.Type)
	})
	smallerFloatFloat := runtime.Signature{
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
		Function: adapter,
		Returns:  runtime.Value{Type: types.Bool},
	}
	smallerFloatInt := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Float,
				Constant: true,
			},
			{
				Name:     "b",
				Type:     types.Integer,
				Constant: true,
			},
		},
		Function: adapter,
		Returns:  runtime.Value{Type: types.Bool},
	}
	smallerIntInt := runtime.Signature{
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
		Function: adapter,
		Returns:  runtime.Value{Type: types.Bool},
	}
	smallerIntFloat := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Integer,
				Constant: true,
			},
			{
				Name:     "b",
				Type:     types.Float,
				Constant: true,
			},
		},
		Function: adapter,
		Returns:  runtime.Value{Type: types.Bool},
	}
	smallerFunction := runtime.Function{
		Signatures: []runtime.Signature{
			smallerIntInt,
			smallerIntFloat,
			smallerFloatFloat,
			smallerFloatInt,
		},
		Source: nil,
	}
	smaller := runtime.Operator{
		Function: smallerFunction,
		Symbol:   "<",
		Constant: true,
	}
	c.Namespace.Store(smaller)
}

func LoadCompare(c *runtime.Context) {
	loadSmaller(c)
}
