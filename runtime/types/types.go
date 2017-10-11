package types

import (
	"fmt"
	"strconv"

	"github.com/tealang/tea-go/runtime"
)

var (
	Any, Bool, Function *runtime.Datatype
	Integer, Float      *runtime.Datatype
	String              *runtime.Datatype
	True, False         runtime.Value
)

func init() {
	Any = &runtime.Datatype{
		Name:   "any",
		Parent: nil,
		Cast: func(v runtime.Value) (runtime.Value, error) {
			return runtime.Value{
				Type: Any,
				Data: v.Data,
			}, nil
		},
		Format: func(v runtime.Value) string {
			v, ok := v.Data.(runtime.Value)
			if !ok {
				return "null"
			}
			return fmt.Sprintf("any<%s>", v.Type.Format(v))
		},
	}
	Integer = &runtime.Datatype{
		Name:   "int",
		Parent: Any,
		Format: func(v runtime.Value) string {
			return fmt.Sprintf("int<%d>", v.Data)
		},
		Cast: func(v runtime.Value) (runtime.Value, error) {
			switch v.Type {
			case nil:
				return runtime.Value{
					Type: Integer,
					Data: int64(0),
				}, nil
			case Integer:
				return v, nil
			case Float:
				return runtime.Value{
					Type: Integer,
					Data: int64(v.Data.(float64)),
				}, nil
			case String:
				i, err := strconv.Atoi(v.Data.(string))
				if err != nil {
					return runtime.Value{}, runtime.ExplicitCastException{From: String, To: Integer}
				}
				return runtime.Value{
					Type: Integer,
					Data: i,
				}, nil
			default:
				return runtime.Value{}, runtime.ExplicitCastException{From: v.Type, To: Integer}
			}
		},
	}
	Float = &runtime.Datatype{
		Name:   "float",
		Parent: Any,
		Format: func(v runtime.Value) string {
			return fmt.Sprintf("float<%f>", v.Data)
		},
	}
	String = &runtime.Datatype{
		Name:   "string",
		Parent: Any,
		Format: func(v runtime.Value) string {
			return fmt.Sprintf("string<'%s'>", v.Data)
		},
		Cast: func(v runtime.Value) (runtime.Value, error) {
			switch v.Type {
			case nil:
				return runtime.Value{
					Type: String,
					Data: "",
				}, nil
			case String:
				return v, nil
			default:
				return runtime.Value{}, runtime.ExplicitCastException{From: v.Type, To: String}
			}
		},
	}
	Function = &runtime.Datatype{
		Name:   "func",
		Parent: Any,
		Cast: func(v runtime.Value) (runtime.Value, error) {
			if v.Type == Function {
				return v, nil
			}
			return runtime.Value{}, runtime.ExplicitCastException{From: v.Type, To: Bool}
		},
		Format: func(v runtime.Value) string {
			return fmt.Sprintf("func<%s>", v.Data)
		},
	}
	Bool = &runtime.Datatype{
		Name:   "bool",
		Parent: Any,
		Cast: func(v runtime.Value) (runtime.Value, error) {
			switch v.Type {
			case Bool:
				return v, nil
			case nil:
				return runtime.Value{
					Type: Bool,
					Data: false,
				}, nil
			default:
				return runtime.Value{}, runtime.ExplicitCastException{From: v.Type, To: Bool}
			}
		},
		Format: func(v runtime.Value) string {
			return fmt.Sprintf("bool<%t>", v.Data)
		},
	}
	True = runtime.Value{
		Type:     Bool,
		Data:     true,
		Constant: true,
		Name:     "true",
	}
	False = runtime.Value{
		Type:     Bool,
		Data:     false,
		Constant: true,
		Name:     "false",
	}
}

func Load(ctx *runtime.Context) {
	ctx.Namespace.Store(Any)
	ctx.Namespace.Store(Function)
	ctx.Namespace.Store(Bool)
	ctx.Namespace.Store(String)
	ctx.Namespace.Store(Integer)
	ctx.Namespace.Store(Float)
}
