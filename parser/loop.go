package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime/nodes"
)

type loopParser struct {
}

func (lp *loopParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	return nil, 0, nil
}
