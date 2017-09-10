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
