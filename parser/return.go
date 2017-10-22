package parser

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/lexer/tokens"
	"github.com/tealang/core/runtime"
	"github.com/tealang/core/runtime/nodes"
)

type returnParser struct {
}

func (returnParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	ctrl := nodes.NewController(runtime.BehaviorReturn)
	if len(input) < 2 {
		return ctrl, 1, nil
	}

	term, n, err := newTermParser().Parse(input[1:])
	if err != nil {
		return ctrl, n + 1, errors.Wrap(err, "parsing return term")
	}
	ctrl.AddBack(term)
	return ctrl, n + 1, nil
}

func newReturnParser() *returnParser {
	return &returnParser{}
}
