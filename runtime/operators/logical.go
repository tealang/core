package operators

import (
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
	"github.com/tealang/tea-go/runtime/types"
)

func loadUnequals(c *runtime.Context) {
	unequalsAny := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Any,
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
				Type: types.Bool,
				Data: a.Data != b.Data,
			}, nil
		}),
		Returns: runtime.Value{Type: types.Bool},
	}
	unequalsFunction := runtime.Function{
		Signatures: []runtime.Signature{
			unequalsAny,
		},
		Source: nil,
	}
	unequals := runtime.Operator{
		Function: unequalsFunction,
		Symbol:   "!=",
		Constant: true,
	}
	c.Namespace.Store(unequals)
}

func loadEquals(c *runtime.Context) {
	equalsAny := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Any,
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
				Type: types.Bool,
				Data: a.Data == b.Data,
			}, nil
		}),
		Returns: runtime.Value{Type: types.Bool},
	}
	equalsFunction := runtime.Function{
		Signatures: []runtime.Signature{
			equalsAny,
		},
		Source: nil,
	}
	equals := runtime.Operator{
		Function: equalsFunction,
		Symbol:   "==",
		Constant: true,
	}
	c.Namespace.Store(equals)
}
func loadNegation(c *runtime.Context) {
	negBool := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Bool,
				Constant: true,
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			var (
				identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
				a         = identA.(runtime.Value)
			)
			return runtime.Value{
				Type: types.Bool,
				Data: !a.Data.(bool),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Bool},
	}
	negFunction := runtime.Function{
		Signatures: []runtime.Signature{
			negBool,
		},
		Source: nil,
	}
	neg := runtime.Operator{
		Function: negFunction,
		Symbol:   "!",
		Constant: true,
	}
	c.Namespace.Store(neg)
}
func loadLogicalXor(c *runtime.Context) {
	xorBool := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Bool,
				Constant: true,
			},
			{
				Name:     "b",
				Type:     types.Bool,
				Constant: true,
			},
		},
		Function: nodes.NewAdapter(func(c *runtime.Context) (runtime.Value, error) {
			var (
				identA, _ = c.Namespace.Find(runtime.SearchIdentifier, "a")
				identB, _ = c.Namespace.Find(runtime.SearchIdentifier, "b")
				a         = identA.(runtime.Value)
				b         = identB.(runtime.Value)
				av        = a.Data.(bool)
				bv        = b.Data.(bool)
			)
			return runtime.Value{
				Type: types.Bool,
				Data: (av || bv) && !(av && bv),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Bool},
	}
	xorFunction := runtime.Function{
		Signatures: []runtime.Signature{
			xorBool,
		},
		Source: nil,
	}
	xor := runtime.Operator{
		Function: xorFunction,
		Symbol:   "^|",
		Constant: true,
	}
	c.Namespace.Store(xor)
}
func loadLogicalOr(c *runtime.Context) {
	orBool := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Bool,
				Constant: true,
			},
			{
				Name:     "b",
				Type:     types.Bool,
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
				Type: types.Bool,
				Data: a.Data.(bool) || b.Data.(bool),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Bool},
	}
	orFunction := runtime.Function{
		Signatures: []runtime.Signature{
			orBool,
		},
		Source: nil,
	}
	or := runtime.Operator{
		Function: orFunction,
		Symbol:   "||",
		Constant: true,
	}
	c.Namespace.Store(or)
}

func loadLogicalAnd(c *runtime.Context) {
	andBool := runtime.Signature{
		Expected: []runtime.Value{
			{
				Name:     "a",
				Type:     types.Bool,
				Constant: true,
			},
			{
				Name:     "b",
				Type:     types.Bool,
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
				Type: types.Bool,
				Data: a.Data.(bool) && b.Data.(bool),
			}, nil
		}),
		Returns: runtime.Value{Type: types.Bool},
	}
	andFunction := runtime.Function{
		Signatures: []runtime.Signature{
			andBool,
		},
		Source: nil,
	}
	and := runtime.Operator{
		Function: andFunction,
		Symbol:   "&&",
		Constant: true,
	}
	c.Namespace.Store(and)
}

func LoadLogical(c *runtime.Context) {
	loadNegation(c)
	loadLogicalAnd(c)
	loadLogicalOr(c)
	loadLogicalXor(c)
	loadEquals(c)
	loadUnequals(c)
}
