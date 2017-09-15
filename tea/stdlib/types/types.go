package types

import "github.com/tealang/tea-go/tea/runtime"

var (
	Any = &runtime.Datatype{
		Name: "any",
	}
	Bool = &runtime.Datatype{
		Name:   "bool",
		Parent: Any,
	}
)

var (
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
)
