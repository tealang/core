package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime/nodes"
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
			return bp.branch, bp.index, err
		}
		// skip offset
		bp.index += n + 1
		stmt, n, err := newSequenceParser(false).Parse(input[bp.index+1:])
		if err != nil {
			return bp.branch, bp.index, err
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
		stmt, n, err := newSequenceParser(true).Parse(input[bp.index+1:])
		if err != nil {
			return bp.branch, bp.index, err
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