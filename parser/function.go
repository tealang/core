package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
)

type functionParser struct {
	active      tokens.Token
	index, size int
	input       []tokens.Token
	literal     bool
	args        []*nodes.Typecast
	returns     *nodes.Typecast
	body        nodes.Node
	alias       string
}

func (fp *functionParser) assignAlias() error {
	if fp.fetch().Type != tokens.Identifier {
		return newUnexpectedTokenException(fp.input[fp.index].Type)
	}
	fp.alias = fp.active.Value
	return nil
}

func (fp *functionParser) fetch() tokens.Token {
	if fp.index >= fp.size {
		return tokens.Token{}
	}
	fp.active = fp.input[fp.index]
	fp.index++
	return fp.active
}

func (fp *functionParser) collectArgs() error {
	if fp.fetch().Type != tokens.LeftParentheses {
		return newUnexpectedTokenException(fp.active.Type)
	}

	var (
		activeArg  nodes.Node
		expectType bool
	)
	for fp.index < fp.size && fp.fetch().Type != tokens.RightParentheses {
		switch fp.active.Type {
		case tokens.Identifier:
			if activeArg != nil {
				if expectType {
					fp.args = append(fp.args, nodes.NewTypecast(fp.active.Value, activeArg))
					expectType = false
				} else {
					return newUnexpectedTokenException(fp.active.Type)
				}
			} else {
				if !expectType {
					activeArg = nodes.NewLiteral(runtime.Value{
						Name: fp.active.Value,
					})
				} else {
					return newUnexpectedTokenException(fp.active.Type)
				}
			}
		case tokens.Operator:
			if fp.active.Value != ":" {
				return newUnexpectedTokenException(fp.active.Type)
			}
			expectType = true
		case tokens.Separator:
			if activeArg == nil {
				return newUnexpectedTokenException(fp.active.Type)
			}
			activeArg = nil
		default:
			return newUnexpectedTokenException(fp.active.Type)
		}
	}
	if fp.active.Type != tokens.RightParentheses {
		return newParseException("Reached unexpected end of program")
	}
	fp.fetch()

	if fp.active.Type == tokens.Operator && fp.active.Value == ":" {
		if fp.fetch().Type != tokens.Identifier {
			return newUnexpectedTokenException(fp.active.Type)
		}
		fp.returns = nodes.NewTypecast(fp.active.Value, nodes.NewLiteral(runtime.Value{}))
		fp.index++
	}

	return nil
}

func (fp *functionParser) collectBody() error {
	stmt, n, err := newSequenceParser(false, 0).Parse(fp.input[fp.index:])
	if err != nil {
		return err
	}
	fp.index += n + 1
	fp.body = stmt
	return nil
}

func (fp *functionParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	fp.index, fp.size = 1, len(input)
	fp.input = input
	if !fp.literal {
		err := fp.assignAlias()
		if err != nil {
			return nil, fp.index, err
		}
	}
	if err := fp.collectArgs(); err != nil {
		return nil, fp.index, err
	}
	if err := fp.collectBody(); err != nil {
		return nil, fp.index, err
	}
	literal := nodes.NewFunctionLiteral(fp.body, fp.returns, fp.args...)
	if fp.literal {
		return literal, fp.index, nil
	}
	return nodes.NewDeclaration(fp.alias, true, literal), fp.index, nil
}

func newFunctionParser(literal bool) *functionParser {
	return &functionParser{literal: literal, args: []*nodes.Typecast{}}
}
