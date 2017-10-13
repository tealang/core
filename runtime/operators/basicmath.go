package operators

import (
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

func loadMultiplication(c *runtime.Context) {
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
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			var (
				identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
				identB, _ = c.Namespace.Find(runtime.SearchIdentifier, "b")
				a         = identA.(runtime.Value)
				b         = identB.(runtime.Value)
			)
			return runtime.Value{
				Type: types.Float,
				Data: a.Data.(float64) * b.Data.(float64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Float},
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
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			var (
				identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
				identB, _ = c.Namespace.Find(runtime.SearchIdentifier, "b")
				a         = identA.(runtime.Value)
				b         = identB.(runtime.Value)
			)
			return runtime.Value{
				Type: types.Integer,
				Data: a.Data.(int64) * b.Data.(int64),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Integer},
	}
	mulFunction := runtime.Function{
		Signatures: []runtime.Signature{
			mulInteger,
			mulFloat,
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
				return runtime.Value{}, runtime.Exception{Message: "Cannot divide by 0"}
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
				return runtime.Value{}, runtime.Exception{Message: "Cannot divide by 0"}
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

func LoadBasicMath(c *runtime.Context) {
	loadAddition(c)
	loadDivision(c)
	loadMultiplication(c)
	loadSubtraction(c)
}
