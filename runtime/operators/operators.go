package operators

import (
	"github.com/tealang/tea-go/runtime"
)

func Load(c *runtime.Context) {
	LoadBasicMath(c)
	LoadLogical(c)
}
