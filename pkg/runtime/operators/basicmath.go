package operators

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/runtime"
	"github.com/tealang/core/pkg/runtime/nodes"
	"github.com/tealang/core/pkg/runtime/types"
)

func loadRemainder(c *runtime.Context) {
	adapter := nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
		var (
			identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
			identB, _ = c.Namespace.Find(runtime.SearchIdentifier, "b")
			a         = identA.(runtime.Value)
			b         = identB.(runtime.Value)
		)
		if b.Data.(int64) == 0 {
			return runtime.Value{}, errors.New("can not divide by 0")
		}
		return runtime.Value{
			Type: types.Integer,
			Data: a.Data.(int64) % b.Data.(int64),
		}, nil
	})
	remInteger := runtime.NewSignature(runtime.Value{Type: types.Integer}, adapter, []runtime.Value{
		{
			Type: types.Integer,
			Name: "a",
		},
		{
			Type: types.Integer,
			Name: "b",
		},
	})
	remFunction := runtime.NewFunction(nil, remInteger)
	remOperator := runtime.Operator{
		Symbol:   "%",
		Function: remFunction,
		Constant: true,
	}
	c.Namespace.Store(remOperator)
}

func loadMultiplication(c *runtime.Context) {
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
				return runtime.Value{}, errors.Wrap(err, "could not multiply")
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "could not multiply")
			}
		}
		switch a.Type {
		case types.Integer:
			return runtime.Value{
				Type: types.Integer,
				Data: a.Data.(int64) * b.Data.(int64),
			}, nil
		case types.Float:
			return runtime.Value{
				Type: types.Float,
				Data: a.Data.(float64) * b.Data.(float64),
			}, nil
		}
		return runtime.Value{}, errors.New("operation * not applicable")
	})
	mulFloat := runtime.Signature{
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
		Returns:  runtime.Value{Type: types.Float},
	}
	mulFloatInteger := runtime.Signature{
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
		Returns:  runtime.Value{Type: types.Float},
	}
	mulInteger := runtime.Signature{
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
		Returns:  runtime.Value{Type: types.Integer},
	}
	mulIntegerFloat := runtime.Signature{
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
		Returns:  runtime.Value{Type: types.Float},
	}
	mulFunction := runtime.Function{
		Signatures: []runtime.Signature{
			mulInteger,
			mulIntegerFloat,
			mulFloat,
			mulFloatInteger,
		},
		Source: nil,
	}
	mul := runtime.Operator{
		Function: mulFunction,
		Symbol:   "*",
		Constant: true,
	}
	c.Namespace.Store(mul)
}
func loadDivision(c *runtime.Context) {
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
				return runtime.Value{}, errors.Wrap(err, "could not divide")
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "could not divide")
			}
		}
		switch a.Type {
		case types.Integer:
			bv := b.Data.(int64)
			if bv == 0 {
				return runtime.Value{}, errors.New("can not divide by 0")
			}
			return runtime.Value{
				Type: types.Integer,
				Data: a.Data.(int64) / bv,
			}, nil
		case types.Float:
			bv := b.Data.(float64)
			if bv == 0 {
				return runtime.Value{}, errors.New("can not divide by 0")
			}
			return runtime.Value{
				Type: types.Float,
				Data: a.Data.(float64) / bv,
			}, nil
		}
		return runtime.Value{}, errors.New("operation / not applicable")
	})
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
		Function: adapter,
		Returns:  runtime.Value{Type: types.Float},
	}
	divFloatInteger := runtime.Signature{
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
		Returns:  runtime.Value{Type: types.Float},
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
		Function: adapter,
		Returns:  runtime.Value{Type: types.Integer},
	}
	divIntegerFloat := runtime.Signature{
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
		Returns:  runtime.Value{Type: types.Float},
	}
	divFunction := runtime.Function{
		Signatures: []runtime.Signature{
			divInteger,
			divIntegerFloat,
			divFloat,
			divFloatInteger,
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
				Type:     types.Any,
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
				Data: a.Data.(string) + b.String(),
			}, nil
		}),
		Returns: runtime.Value{Type: types.String},
	}
	addAdapter := nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
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
				return runtime.Value{}, errors.Wrap(err, "could not add")
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "could not add")
			}
		}
		switch a.Type {
		case types.Integer:
			return runtime.Value{
				Type: types.Integer,
				Data: a.Data.(int64) + b.Data.(int64),
			}, nil
		case types.Float:
			return runtime.Value{
				Type: types.Float,
				Data: a.Data.(float64) + b.Data.(float64),
			}, nil
		}
		return runtime.Value{}, errors.New("operation + not applicable")
	})
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
		Function: addAdapter,
		Returns:  runtime.Value{Type: types.Float},
	}
	addFloatInteger := runtime.Signature{
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
		Function: addAdapter,
		Returns:  runtime.Value{Type: types.Float},
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
		Function: addAdapter,
		Returns:  runtime.Value{Type: types.Integer},
	}
	addIntegerFloat := runtime.Signature{
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
		Function: addAdapter,
		Returns:  runtime.Value{Type: types.Float},
	}
	addFunction := runtime.Function{
		Signatures: []runtime.Signature{
			joinString,
			addInteger,
			addIntegerFloat,
			addFloat,
			addFloatInteger,
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
	subAdapter := nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
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
				return runtime.Value{}, errors.Wrap(err, "could not subtract")
			}
		} else if b.Type == types.Float {
			a, err = b.Type.Cast(b)
			if err != nil {
				return runtime.Value{}, errors.Wrap(err, "could not subtract")
			}
		}
		switch a.Type {
		case types.Integer:
			return runtime.Value{
				Type: types.Integer,
				Data: a.Data.(int64) - b.Data.(int64),
			}, nil
		case types.Float:
			return runtime.Value{
				Type: types.Float,
				Data: a.Data.(float64) - b.Data.(float64),
			}, nil
		}
		return runtime.Value{}, errors.New("operator - not applicable")
	})
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
		Function: subAdapter,
		Returns:  runtime.Value{Type: types.Float},
	}
	subFloatInteger := runtime.Signature{
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
		Function: subAdapter,
		Returns:  runtime.Value{Type: types.Float},
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
		Function: subAdapter,
		Returns:  runtime.Value{Type: types.Integer},
	}
	subIntegerFloat := runtime.Signature{
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
		Function: subAdapter,
		Returns:  runtime.Value{Type: types.Float},
	}
	subFunction := runtime.Function{
		Signatures: []runtime.Signature{
			subInteger,
			subIntegerFloat,
			subFloat,
			subFloatInteger,
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

// LoadBasicMath loads basic math operators like addition and subtraction into the runtime context.
func LoadBasicMath(c *runtime.Context) {
	loadAddition(c)
	loadDivision(c)
	loadMultiplication(c)
	loadSubtraction(c)
	loadRemainder(c)
}
