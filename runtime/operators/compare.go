package operators

import (
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

func loadGreater(c *runtime.Context) {
	adapter := nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
		var (
			identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
			identB, _ = c.Namespace.Find(runtime.SearchIdentifier, "b")
			a         = identA.(runtime.Value)
			b         = identB.(runtime.Value)
			err       error
		)
		if a.Type == types.Float {
			b, err = a.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, err
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, err
			}
		}
		switch a.Type {
		case types.Integer:
			return runtime.Value{
				Type: types.Bool,
				Data: a.Data.(int64) > b.Data.(int64),
			}, nil
		case types.Float:
			return runtime.Value{
				Type: types.Bool,
				Data: a.Data.(float64) > b.Data.(float64),
			}, nil
		}
		return runtime.Value{}, runtime.NewNotApplicableException("<", a.Type, b.Type)
	})
	greaterFloat := runtime.Signature{
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
	greaterFloatInt := runtime.Signature{
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
	greaterInt := runtime.Signature{
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
	greaterIntFloat := runtime.Signature{
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
	greaterFunction := runtime.Function{
		Signatures: []runtime.Signature{
			greaterFloat,
			greaterFloatInt,
			greaterInt,
			greaterIntFloat,
		},
		Source: nil,
	}
	greater := runtime.Operator{
		Function: greaterFunction,
		Symbol:   ">",
		Constant: true,
	}
	c.Namespace.Store(greater)
}

func loadSmaller(c *runtime.Context) {
	adapter := nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
		var (
			identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
			identB, _ = c.Namespace.Find(runtime.SearchIdentifier, "b")
			a         = identA.(runtime.Value)
			b         = identB.(runtime.Value)
			err       error
		)
		if a.Type == types.Float {
			b, err = a.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, err
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, err
			}
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
	loadGreater(c)
}
