package parser

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/lexer/tokens"
	"github.com/tealang/core/pkg/runtime"
	"github.com/tealang/core/pkg/runtime/nodes"
)

type parameterizedSequenceParser struct {
	active      tokens.Token
	index, size int
	input       []tokens.Token
	args        []*nodes.Type
	returns     *nodes.Type
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
		return errors.Errorf("did expect left parentheses, got %s", sp.active.Type)
	}

	var (
		activeArgs []nodes.Node
		expectType bool
	)
	for sp.index < sp.size && sp.fetch().Type != tokens.RightParentheses {
		switch sp.active.Type {
		case tokens.Identifier:
			if activeArgs != nil {
				if expectType {
					typenode, offset, err := newTypeParser().Parse(sp.input[sp.index-1:])
					if err != nil {
						return errors.Wrap(err, "failed to parse param")
					}
					sp.index += offset - 1
					for _, arg := range activeArgs {
						sp.args = append(sp.args, nodes.NewType(typenode.(*nodes.Type).Tree, arg))
					}
					activeArgs = nil
					expectType = false
				} else {
					activeArgs = append(activeArgs, nodes.NewLiteral(runtime.Value{Name: sp.active.Value}))
				}
			} else {
				if !expectType {
					activeArgs = []nodes.Node{
						nodes.NewLiteral(runtime.Value{Name: sp.active.Value}),
					}
				} else {
					return errors.New("did not expect identifier")
				}
			}
		case tokens.Operator:
			if sp.active.Value != ":" {
				return errors.Errorf("expected typecast operator, got %s", sp.active.Value)
			}
			expectType = true
		case tokens.Separator:
		default:
			return errors.Errorf("did not expect token %s", sp.active.Type)
		}
	}
	if sp.active.Type != tokens.RightParentheses {
		return errors.New("expected right parentheses, reached unexpected end of program")
	}
	sp.fetch()

	if sp.active.Type == tokens.Operator && sp.active.Value == ":" {
		if sp.fetch().Type != tokens.Identifier {
			return errors.Errorf("expected results cast, got %s", sp.active.Type)
		}
		typenode, offset, err := newTypeParser().Parse(sp.input[sp.index-1:])
		if err != nil {
			return errors.Wrap(err, "failed to parse return type")
		}
		sp.index += offset - 1
		sp.returns = typenode.(*nodes.Type)
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

func (sp *parameterizedSequenceParser) Parse(input []tokens.Token) ([]*nodes.Type, nodes.Node, *nodes.Type, int, error) {
	sp.index, sp.size = 0, len(input)
	sp.input = input
	if err := sp.collectArgs(); err != nil {
		return nil, nil, nil, sp.index, errors.Wrap(err, "failed collecting args")
	}
	if err := sp.collectBody(); err != nil {
		return nil, nil, nil, sp.index, errors.Wrap(err, "failed parsing body")
	}
	return sp.args, sp.body, sp.returns, sp.index, nil
}

func newParameterizedSequenceParser() *parameterizedSequenceParser {
	return &parameterizedSequenceParser{args: make([]*nodes.Type, 0)}
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
		return errors.Errorf("expected function alias, got %s", fp.active.Type)
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
			return nil, fp.index, errors.Wrap(err, "failed to parse function")
		}
	}
	args, body, returns, n, err := newParameterizedSequenceParser().Parse(input[fp.index:])
	if err != nil {
		return nil, fp.index, errors.Wrap(err, "failed to parse function")
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
