package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime/nodes"
)

type loopParser struct {
	index, size int
}

func (lp *loopParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	lp.index, lp.size = 1, len(input)

	entry, n, err := newTermParser().Parse(input[lp.index:])
	if err != nil {
		return nil, lp.index, err
	}
	// check if single-tier loop head
	if input[lp.index+n].Type == tokens.LeftBlock {
		lp.index += n
		body, n, err := newSequenceParser(false, 0).Parse(input[lp.index+1:])
		if err != nil {
			return nil, lp.index, err
		}
		// ignore right block
		lp.index += n + 2
		return nodes.NewLoop(entry, body), lp.index, nil
	}

	// handle three-tier loop
	sequ, n, err := newSequenceParser(false, 3).Parse(input[lp.index:])
	if err != nil {
		return nil, lp.index, err
	}
	head := sequ.(*nodes.Sequence)
	lp.index += n

	if input[lp.index].Type != tokens.LeftBlock {
		return nil, lp.index, newUnexpectedTokenException(input[lp.index].Type)
	}
	if len(head.Childs) != 3 {
		return nil, lp.index, newParseException("Expected three-tier for-loop")
	}

	body, n, err := newSequenceParser(false, 0).Parse(input[lp.index+1:])
	if err != nil {
		return nil, lp.index, err
	}
	lp.index += n + 2

	return nodes.NewSequence(true, head.Childs[0], nodes.NewLoop(head.Childs[1], nodes.NewSequence(false, body, head.Childs[2]))), lp.index, nil
}

func newLoopParser() *loopParser {
	return &loopParser{}
}
