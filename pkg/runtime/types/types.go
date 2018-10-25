package types

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/runtime"
)

// Default types defined by runtime.
var (
	Any, Bool, Function *runtime.Datatype
	Integer, Float      *runtime.Datatype
	String              *runtime.Datatype
	Array               *runtime.Datatype
)

// Boolean values.
var (
	True, False runtime.Value
)

func init() {
	Any = &runtime.Datatype{
		Name:   "any",
		Parent: nil,
		Cast: func(v runtime.Value, f []runtime.Typeflag) (runtime.Value, error) {
			switch v.Type {
			case nil:
				return runtime.Value{
					Typeflag: runtime.T(Any),
					Name:     v.Name,
					Data:     nil,
				}, nil
			case Any:
				return runtime.Value{
					Typeflag: v.Typeflag,
					Data:     v.Data,
					Name:     v.Name,
				}, nil
			default:
				return runtime.Value{
					Typeflag: runtime.Typeflag{
						Type:   Any,
						Params: []runtime.Typeflag{v.Typeflag},
					},
					Name: v.Name,
					Data: v.Data,
				}, nil
			}
		},
		Format: func(v runtime.Value) string {
			if v.Data != nil {
				return fmt.Sprint(v.Data)
			}
			return "nil"
		},
	}
	Array = &runtime.Datatype{
		Name:   "array",
		Parent: Any,
		Cast: func(v runtime.Value, f []runtime.Typeflag) (runtime.Value, error) {
			switch v.Type {
			case nil:
				return runtime.Value{
					Typeflag: v.Typeflag,
					Data: nil,
					Name: v.Name,
				}, nil
			case Array:
				return runtime.Value{
					Typeflag: v.Typeflag,
					Data:     v.Data,
					Name:     v.Name,
				}, nil
			default:
				return runtime.Value{}, errors.Errorf("can not cast %s to array", v.Type)
			}
		},
		Format: func(v runtime.Value) string {
			return "array"
		},
	}
	Integer = &runtime.Datatype{
		Name:   "int",
		Parent: Any,
		Format: func(v runtime.Value) string {
			return fmt.Sprint(v.Data)
		},
		Cast: func(v runtime.Value, f []runtime.Typeflag) (runtime.Value, error) {
			if len(f) != 0 {
				return runtime.Value{}, errors.New("unsupported type parameters")
			}
			switch v.Type {
			case nil:
				return runtime.Value{
					Typeflag: runtime.T(Integer),
					Data:     int64(0),
					Name:     v.Name,
				}, nil
			case Integer:
				return v, nil
			case Float:
				return runtime.Value{
					Typeflag: runtime.T(Integer),
					Data:     int64(v.Data.(float64)),
					Name:     v.Name,
				}, nil
			case String:
				i, err := strconv.Atoi(v.Data.(string))
				if err != nil {
					return runtime.Value{}, errors.Wrap(err, "can not cast string to int")
				}
				return runtime.Value{
					Typeflag: runtime.T(Integer),
					Data:     i,
					Name:     v.Name,
				}, nil
			default:
				return runtime.Value{}, errors.Errorf("can not cast %s to int", v.Type)
			}
		},
	}
	Float = &runtime.Datatype{
		Name:   "float",
		Parent: Any,
		Format: func(v runtime.Value) string {
			return fmt.Sprint(v.Data)
		},
		Cast: func(v runtime.Value, f []runtime.Typeflag) (runtime.Value, error) {
			if len(f) != 0 {
				return runtime.Value{}, errors.New("unsupported type parameters")
			}
			switch v.Type {
			case nil:
				return runtime.Value{
					Typeflag: runtime.T(Float),
					Data:     float64(0),
					Name:     v.Name,
				}, nil
			case Integer:
				return runtime.Value{
					Typeflag: runtime.T(Float),
					Data:     float64(v.Data.(int64)),
					Name:     v.Name,
				}, nil
			case Float:
				return v, nil
			default:
				return runtime.Value{}, errors.Errorf("can not cast %s to float", v.Type)
			}
		},
	}
	String = &runtime.Datatype{
		Name:   "string",
		Parent: Any,
		Format: func(v runtime.Value) string {
			return fmt.Sprintf(`%s`, v.Data)
		},
		Cast: func(v runtime.Value, f []runtime.Typeflag) (runtime.Value, error) {
			if len(f) != 0 {
				return runtime.Value{}, errors.New("unsupported type parameters")
			}
			switch v.Type {
			case nil:
				return runtime.Value{
					Typeflag: runtime.T(String),
					Data:     "",
					Name:     v.Name,
				}, nil
			case String:
				return v, nil
			default:
				return runtime.Value{}, errors.Errorf("can not cast %s to string", v.Type)
			}
		},
	}
	Function = &runtime.Datatype{
		Name:   "func",
		Parent: Any,
		Cast: func(v runtime.Value, f []runtime.Typeflag) (runtime.Value, error) {
			if len(f) != 0 {
				return runtime.Value{}, errors.New("unsupported type parameters")
			}
			if v.Type == Function {
				return v, nil
			}
			return runtime.Value{}, errors.Errorf("can not cast %s to func", v.Type)
		},
		Format: func(v runtime.Value) string {
			return fmt.Sprintf("func<%s>", v.Data)
		},
	}
	Bool = &runtime.Datatype{
		Name:   "bool",
		Parent: Any,
		Cast: func(v runtime.Value, f []runtime.Typeflag) (runtime.Value, error) {
			if len(f) != 0 {
				return runtime.Value{}, errors.New("unsupported type parameters")
			}
			switch v.Type {
			case Bool:
				return v, nil
			case nil:
				return runtime.Value{
					Typeflag: runtime.T(Bool),
					Name:     v.Name,
					Data:     false,
				}, nil
			default:
				return runtime.Value{}, errors.Errorf("can not cast %s to bool", v.Type)
			}
		},
		Format: func(v runtime.Value) string {
			return fmt.Sprint(v.Data)
		},
	}
	True = runtime.Value{
		Typeflag: runtime.T(Bool),
		Data:     true,
		Constant: true,
		Name:     "true",
	}
	False = runtime.Value{
		Typeflag: runtime.T(Bool),
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
	ctx.Namespace.Store(Array)
}
