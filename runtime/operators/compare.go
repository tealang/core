package operators

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/runtime"
	"github.com/tealang/core/runtime/nodes"
	"github.com/tealang/core/runtime/types"
)

func loadGreaterEqual(c *runtime.Context) {
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
				return runtime.Value{}, errors.Wrap(err, "can not compare")
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "can not compare")
			}
		}
		switch a.Type {
		case types.Integer:
			return runtime.Value{
				Type: types.Bool,
				Data: a.Data.(int64) >= b.Data.(int64),
			}, nil
		case types.Float:
			return runtime.Value{
				Type: types.Bool,
				Data: a.Data.(float64) >= b.Data.(float64),
			}, nil
		}
		return runtime.Value{}, errors.New("operator >= not applicable")
	})
	greaterEqualFloat := runtime.Signature{
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
	greaterEqualFloatInt := runtime.Signature{
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
	greateEqualInt := runtime.Signature{
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
	greaterEqualIntFloat := runtime.Signature{
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
	greaterEqualFunction := runtime.Function{
		Signatures: []runtime.Signature{
			greaterEqualFloat,
			greaterEqualFloatInt,
			greateEqualInt,
			greaterEqualIntFloat,
		},
		Source: nil,
	}
	greaterEqual := runtime.Operator{
		Function: greaterEqualFunction,
		Symbol:   ">=",
		Constant: true,
	}
	c.Namespace.Store(greaterEqual)
}

func loadSmallerEqual(c *runtime.Context) {
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
				return runtime.Value{}, errors.Wrap(err, "can not compare")
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "can not compare")
			}
		}
		switch a.Type {
		case types.Integer:
			return runtime.Value{
				Type: types.Bool,
				Data: a.Data.(int64) <= b.Data.(int64),
			}, nil
		case types.Float:
			return runtime.Value{
				Type: types.Bool,
				Data: a.Data.(float64) <= b.Data.(float64),
			}, nil
		}
		return runtime.Value{}, errors.New("operator <= not applicable")
	})
	smallerEqualFloat := runtime.Signature{
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
	smallerEqualFloatInt := runtime.Signature{
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
	smallerEqualInt := runtime.Signature{
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
	smallerEqualIntFloat := runtime.Signature{
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
	smallerEqualFunction := runtime.Function{
		Signatures: []runtime.Signature{
			smallerEqualInt,
			smallerEqualIntFloat,
			smallerEqualFloat,
			smallerEqualFloatInt,
		},
		Source: nil,
	}
	smallerEqual := runtime.Operator{
		Function: smallerEqualFunction,
		Symbol:   "<=",
		Constant: true,
	}
	c.Namespace.Store(smallerEqual)
}

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
				return runtime.Value{}, errors.Wrap(err, "can not compare")
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "can not compare")
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
		return runtime.Value{}, errors.New("operator > not applicable")
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
				return runtime.Value{}, errors.Wrap(err, "can not compare")
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "can not compare")
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
		return runtime.Value{}, errors.New("operator < not applicable")
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

// LoadCompare loads comparison operators into the runtime context.
func LoadCompare(c *runtime.Context) {
	loadSmaller(c)
	loadSmallerEqual(c)
	loadGreater(c)
	loadGreaterEqual(c)
}
