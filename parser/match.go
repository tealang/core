package parser

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/lexer/tokens"
	"github.com/tealang/core/runtime/nodes"
)

type matchParser struct {
	index, size int
	cases       []nodes.Node
}

func (mp *matchParser) parseCase(input []tokens.Token, isDefault bool) error {
	var matchTo nodes.Node
	if !isDefault {
		term, offset, err := newTermParser().Parse(input[mp.index:])
		if err != nil {
			return errors.Wrap(err, "failed to build case")
		}
		mp.index += offset
		matchTo = term
	}
	if input[mp.index].Type != tokens.LeftBlock {
		return errors.Errorf("expected left block")
	}
	mp.index++
	body, offset, err := newSequenceParser(false, 0).Parse(input[mp.index:])
	if err != nil {
		return errors.Wrap(err, "failed to build case body")
	}
	mp.index += offset
	if input[mp.index].Type != tokens.RightBlock {
		return errors.Errorf("expected right block")
	}
	mp.index++
	if !isDefault {
		mp.cases = append(mp.cases, nodes.NewCase(matchTo, body))
	} else {
		mp.cases = append(mp.cases, body)
	}
	return nil
}

func (mp *matchParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	mp.index, mp.size = 0, len(input)

	// Match ...
	if input[mp.index].Type != tokens.Identifier || input[mp.index].Value != matchKeyword {
		return nil, mp.index, errors.Errorf("expected match keyword")
	}
	mp.index++

	// match <term>
	term, offset, err := newTermParser().Parse(input[mp.index:])
	if err != nil {
		return nil, mp.index, errors.Wrap(err, "failed to build header")
	}
	mp.index += offset

	// match <term> {
	if input[mp.index].Type != tokens.LeftBlock {
		return nil, mp.index, errors.Errorf("failed to build match: expected left block")
	}
	mp.index++

	mp.cases = nil
	for mp.index < mp.size && input[mp.index].Type == tokens.Identifier && input[mp.index].Value == caseKeyword {
		mp.index++
		if err := mp.parseCase(input, false); err != nil {
			return nil, mp.index, errors.Errorf("failed to build match: %v", err)
		}
	}

	if mp.index < mp.size && input[mp.index].Type == tokens.Identifier && input[mp.index].Value == defaultKeyword {
		mp.index++
		if err := mp.parseCase(input, true); err != nil {
			return nil, mp.index, errors.Errorf("failed to build match default: %v")
		}
	}

	if mp.index >= mp.size || input[mp.index].Type != tokens.RightBlock {
		return nil, mp.index, errors.Errorf("expected right block")
	}
	mp.index++

	return nodes.NewMatch(term, mp.cases...), mp.index, nil
}

func newMatchParser() *matchParser {
	return &matchParser{}
}
