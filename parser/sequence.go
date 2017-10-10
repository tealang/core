package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
)

func newSequenceParser(substitute bool) *sequenceParser {
	return &sequenceParser{substitute}
}

type sequenceParser struct {
	Substitute bool
}

func (sp *sequenceParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	var (
		index = 0
		seq   = nodes.NewSequence(sp.Substitute)
	)
	for ; index < len(input); index++ {
		requireEndStatement := true
		active := input[index]
		switch active.Type {
		case tokens.RightBlock:
			return seq, index, nil
		case tokens.LeftBlock:
			item, n, err := newSequenceParser(true).Parse(input[index+1:])
			if err != nil {
				return seq, index, err
			}
			index += n
			seq.AddBack(item)
		case tokens.Identifier:
			switch active.Value {
			case "let", "var":
				stmt, n, err := newDeclarationParser().Parse(input[index:])
				if err != nil {
					return seq, index, err
				}
				seq.AddBack(stmt)
				index += n
			case "return":
				stmt, n, err := newReturnParser().Parse(input[index:])
				if err != nil {
					return seq, index, err
				}
				seq.AddBack(stmt)
				index += n
			case "break":
				seq.AddBack(nodes.NewController(runtime.BehaviorBreak))
				index++
			case "continue":
				seq.AddBack(nodes.NewController(runtime.BehaviorContinue))
				index++
			default:
				term, n, err := newTermParser().Parse(input[index:])
				if err != nil {
					return seq, index, err
				}
				seq.AddBack(term)
				index += n
			}
		default:
			term, n, err := newTermParser().Parse(input[index:])
			if err != nil {
				return seq, index, err
			}
			seq.AddBack(term)
			index += n
		}

		if index >= len(input) {
			requireEndStatement = false
		}

		if requireEndStatement && input[index].Type != tokens.Statement {
			return seq, index, ParseException{"Expected end statement"}
		}
	}
	return seq, index, nil
}
