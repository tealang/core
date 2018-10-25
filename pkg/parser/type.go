package parser

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/lexer/tokens"
	"github.com/tealang/core/pkg/runtime/nodes"
)

type typeParser struct {
	index, size int
	input       []tokens.Token
	active, next      tokens.Token
}

func (tp *typeParser) fetch() tokens.Token {
	if tp.index >= tp.size {
		return tokens.Token{}
	}
	tp.active = tp.input[tp.index]
	tp.index++
	if tp.index < tp.size {
		tp.next = tp.input[tp.index]
	} else {
		tp.next = tokens.Token{}
	}
	return tp.active
}

func (tp *typeParser) tree() (nodes.Typetree, error) {
	if tp.fetch().Type != tokens.Identifier {
		return nodes.Typetree{}, errors.New("missing typename")
	}
	tree := nodes.Typetree{
		Name: tp.active.Value,
	}
	if tp.next.Type != tokens.Operator || tp.next.Value != "<" {
		return tree, nil
	}
	tp.fetch()
	for tp.next.Type == tokens.Identifier {
		subtree, err := tp.tree()
		if err != nil {
			return tree, err
		}
		tree.Params = append(tree.Params, subtree)
		if tp.fetch().Type == tokens.Operator && tp.active.Value == ">" {
			return tree, nil
		}
		if tp.active.Type != tokens.Separator {
			return tree, errors.New("expected separator")
		}
	}
	return tree, nil
}

func (tp *typeParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	tp.index, tp.size = 0, len(input)
	tp.input = input

	tree, err := tp.tree()
	if err != nil {
		return nil, tp.index, errors.Wrap(err, "could not parse type")
	}
	return nodes.NewType(tree), tp.index, nil
}

func newTypeParser() *typeParser {
	return &typeParser{}
}
