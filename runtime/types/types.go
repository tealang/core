package types

import (
	"fmt"

	"github.com/tealang/tea-go/runtime"
)

var (
	Any, Bool, Function    *runtime.Datatype
	Number, Integer, Float *runtime.Datatype
	String                 *runtime.Datatype
	True, False            runtime.Value
)

func init() {
	Any = &runtime.Datatype{
		Name:   "any",
		Parent: nil,
		Cast: func(v runtime.Value) (runtime.Value, error) {
			return runtime.Value{
				Type: Any,
				Data: v,
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
	Number = &runtime.Datatype{
		Name:   "number",
		Parent: Any,
	}
	Integer = &runtime.Datatype{
		Name:   "int",
		Parent: Number,
		Format: func(v runtime.Value) string {
			return fmt.Sprintf("int<%d>", v.Data)
		},
	}
	Float = &runtime.Datatype{
		Name:   "float",
		Parent: Number,
		Format: func(v runtime.Value) string {
			return fmt.Sprintf("float<%f>", v.Data)
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
			if v.Type == Bool {
				return v, nil
			}
			return runtime.Value{}, runtime.ExplicitCastException{From: v.Type, To: Bool}
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
