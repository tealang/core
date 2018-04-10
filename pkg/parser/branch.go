package parser

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/lexer/tokens"
	"github.com/tealang/core/pkg/runtime/nodes"
)

type branchParser struct {
	index, size int
	branch      *nodes.Branch
}

func (bp *branchParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	bp.index, bp.size = 0, len(input)
	for input[bp.index].Value == ifKeyword {
		condition, n, err := newTermParser().Parse(input[bp.index+1:])
		if err != nil {
			return bp.branch, bp.index, errors.Wrap(err, "can not parse condition")
		}
		// skip offset
		bp.index += n + 1
		stmt, n, err := newSequenceParser(false, 0).Parse(input[bp.index+1:])
		if err != nil {
			return bp.branch, bp.index, errors.Wrap(err, "can not parse body")
		}
		// skip right block and offset
		bp.index += n + 2
		bp.branch.AddBack(nodes.NewConditional(condition, stmt))

		if bp.index >= bp.size || input[bp.index].Value != elseKeyword {
			return bp.branch, bp.index, nil
		}

		bp.index++
	}

	// add else if needed
	if input[bp.index].Type == tokens.LeftBlock {
		// must substitute, not encapsulated in conditional
		stmt, n, err := newSequenceParser(true, 0).Parse(input[bp.index+1:])
		if err != nil {
			return bp.branch, bp.index, errors.Wrap(err, "can not parse else body")
		}
		bp.index += n + 2
		bp.branch.AddBack(stmt)
	}

	return bp.branch, bp.index, nil
}

func newBranchParser() *branchParser {
	return &branchParser{
		branch: nodes.NewBranch(),
	}
}
