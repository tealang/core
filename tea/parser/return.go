package parser

import (
	"github.com/tealang/tea-go/tea/lexer/tokens"
	"github.com/tealang/tea-go/tea/runtime"
	"github.com/tealang/tea-go/tea/runtime/nodes"
)

func GenerateReturn(input []tokens.Token) (nodes.Node, int, error) {
	ctrl := nodes.NewController(runtime.BehaviorReturn)
	if len(input) < 2 {
		return ctrl, 1, nil
	}

	term, n, err := GenerateTerm(input[1:])
	if err != nil {
		return ctrl, n + 1, err
	}
	ctrl.AddBack(term)
	return ctrl, n + 1, nil
}
