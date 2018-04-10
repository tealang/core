package operators

import (
	"github.com/tealang/core/pkg/runtime"
)

func Load(c *runtime.Context) {
	LoadBasicMath(c)
	LoadLogical(c)
	LoadCompare(c)
}
