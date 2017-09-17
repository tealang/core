package parser

import (
	"github.com/tealang/tea-go/tea/lexer/tokens"
	"github.com/tealang/tea-go/tea/runtime"
	"github.com/tealang/tea-go/tea/runtime/nodes"
)

func GenerateDeclaration(input []tokens.Token) (nodes.Node, int, error) {
	var (
		index = 0
		decl  = nodes.NewMultiDeclaration([]string{}, false)
	)
	switch input[index].Value {
	case "var":
		decl.Constant = false
	case "let":
		decl.Constant = true
	}
	index++

	// collect aliases
	var (
		casts               = make([]string, 0)
		expectTypeInfo      bool
		expectAssignment    bool
		stopCollectingAlias bool
	)
	for !stopCollectingAlias && index < len(input) {
		active := input[index]
		switch active.Type {
		case tokens.Identifier:
			if expectTypeInfo {
				for len(casts) < len(decl.Alias) {
					casts = append(casts, active.Value)
				}
			} else {
				decl.Alias = append(decl.Alias, active.Value)
			}
			expectTypeInfo = false
		case tokens.Statement:
			stopCollectingAlias = true
		case tokens.Separator:
		case tokens.Operator:
			if active.Value == ":" {
				expectTypeInfo = true
			} else if active.Value == "=" {
				expectAssignment = true
				stopCollectingAlias = true
			} else {
				return decl, index, ParseException{"Unexpected operator"}
			}
		default:
			return decl, index, ParseException{"Unexpected token"}
		}
		if !stopCollectingAlias {
			index++
		}
	}

	// handle if there is no direct assignment
	if !expectAssignment {
		if len(casts) != len(decl.Alias) {
			return decl, index, ParseException{"Required type information not found"}
		}
		for _, t := range casts {
			decl.AddBack(nodes.NewTypecast(t, nodes.NewLiteral(runtime.Value{})))
		}
		return decl, index, nil
	}

	// collect values
	termIteration := 0
	for ; index < len(input); index++ {
		term, n, err := GenerateTerm(input[index:])
		if err != nil {
			return decl, index, err
		}
		index += n

		if termIteration < len(decl.Alias) {
			term = nodes.NewTypecast(casts[termIteration])
			termIteration++
		}
		decl.AddBack(term)
	}

	return decl, index, nil
}
