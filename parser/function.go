package parser

import (
	"github.com/tealang/core/lexer/tokens"
	"github.com/tealang/core/runtime"
	"github.com/tealang/core/runtime/nodes"
)

type parameterizedSequenceParser struct {
	active      tokens.Token
	index, size int
	input       []tokens.Token
	args        []*nodes.Typecast
	returns     *nodes.Typecast
	body        nodes.Node
}

func (sp *parameterizedSequenceParser) fetch() tokens.Token {
	if sp.index >= sp.size {
		return tokens.Token{}
	}
	sp.active = sp.input[sp.index]
	sp.index++
	return sp.active
}

func (sp *parameterizedSequenceParser) collectArgs() error {
	if sp.fetch().Type != tokens.LeftParentheses {
		return newUnexpectedTokenException(sp.active.Type)
	}

	var (
		activeArg  nodes.Node
		expectType bool
	)
	for sp.index < sp.size && sp.fetch().Type != tokens.RightParentheses {
		switch sp.active.Type {
		case tokens.Identifier:
			if activeArg != nil {
				if expectType {
					sp.args = append(sp.args, nodes.NewTypecast(sp.active.Value, activeArg))
					expectType = false
				} else {
					return newUnexpectedTokenException(sp.active.Type)
				}
			} else {
				if !expectType {
					activeArg = nodes.NewLiteral(runtime.Value{
						Name: sp.active.Value,
					})
				} else {
					return newUnexpectedTokenException(sp.active.Type)
				}
			}
		case tokens.Operator:
			if sp.active.Value != ":" {
				return newUnexpectedTokenException(sp.active.Type)
			}
			expectType = true
		case tokens.Separator:
			if activeArg == nil {
				return newUnexpectedTokenException(sp.active.Type)
			}
			activeArg = nil
		default:
			return newUnexpectedTokenException(sp.active.Type)
		}
	}
	if sp.active.Type != tokens.RightParentheses {
		return newParseException("Reached unexpected end of program")
	}
	sp.fetch()

	if sp.active.Type == tokens.Operator && sp.active.Value == ":" {
		if sp.fetch().Type != tokens.Identifier {
			return newUnexpectedTokenException(sp.active.Type)
		}
		sp.returns = nodes.NewTypecast(sp.active.Value, nodes.NewLiteral(runtime.Value{}))
		sp.index++
	}

	return nil
}

func (sp *parameterizedSequenceParser) collectBody() error {
	stmt, n, err := newSequenceParser(false, 0).Parse(sp.input[sp.index:])
	if err != nil {
		return err
	}
	sp.index += n + 1
	sp.body = stmt
	return nil
}
func (sp *parameterizedSequenceParser) Parse(input []tokens.Token) ([]*nodes.Typecast, nodes.Node, *nodes.Typecast, int, error) {
	sp.index, sp.size = 0, len(input)
	sp.input = input
	if err := sp.collectArgs(); err != nil {
		return nil, nil, nil, sp.index, err
	}
	if err := sp.collectBody(); err != nil {
		return nil, nil, nil, sp.index, err
	}
	return sp.args, sp.body, sp.returns, sp.index, nil
}

func newParameterizedSequenceParser() *parameterizedSequenceParser {
	return &parameterizedSequenceParser{args: make([]*nodes.Typecast, 0)}
}

type functionParser struct {
	active      tokens.Token
	index, size int
	input       []tokens.Token
	literal     bool
	alias       string
}

func (fp *functionParser) assignAlias() error {
	if fp.fetch().Type != tokens.Identifier {
		return newParseException("Expected function name")
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

func (fp *functionParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	fp.index, fp.size = 1, len(input)
	fp.input = input
	if !fp.literal {
		err := fp.assignAlias()
		if err != nil {
			return nil, fp.index, err
		}
	}
	args, body, returns, n, err := newParameterizedSequenceParser().Parse(input[fp.index:])
	if err != nil {
		return nil, fp.index, err
	}
	fp.index += n
	literal := nodes.NewFunctionLiteral(body, returns, args...)
	if fp.literal {
		return literal, fp.index, nil
	}
	return nodes.NewDeclaration(fp.alias, true, literal), fp.index, nil
}

func newFunctionParser(literal bool) *functionParser {
	return &functionParser{literal: literal}
}
