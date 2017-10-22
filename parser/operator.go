package parser

import (
	"github.com/tealang/core/lexer/tokens"
	"github.com/tealang/core/runtime/nodes"
)

type operatorParser struct {
	index, size int
	input       []tokens.Token
	symbol      string
	active      tokens.Token
}

func (op *operatorParser) fetch() tokens.Token {
	if op.index >= op.size {
		return tokens.Token{}
	}
	op.active = op.input[op.index]
	op.index++
	return op.active
}

func (op *operatorParser) assignSymbol() error {
	if op.fetch().Type != tokens.Operator {
		return newUnexpectedTokenException(op.input[op.index])
	}
	op.symbol = op.active.Value
	return nil
}

func (op *operatorParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	op.index, op.size = 0, len(input)
	op.input = input

	if op.fetch().Type != tokens.Identifier && op.active.Value != operatorKeyword {
		return nil, op.index, newParseException("Expected operator keyword")
	}
	if err := op.assignSymbol(); err != nil {
		return nil, op.index, err
	}
	args, body, returns, n, err := newParameterizedSequenceParser().Parse(input[op.index:])
	if err != nil {
		return nil, op.index, err
	}
	op.index += n
	return nodes.NewOperatorDefinition(op.symbol, body, returns, args...), op.index, nil
}

func newOperatorParser() *operatorParser {
	return &operatorParser{}
}
