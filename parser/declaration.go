package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
)

type declarationParser struct {
	ExpectTypeInformation bool
	ExpectAssignment      bool
	Casts                 []string
	Declaration           *nodes.Declaration
	Index                 int
}

func newDeclarationParser() *declarationParser {
	return &declarationParser{
		Casts:       make([]string, 0),
		Declaration: nodes.NewMultiDeclaration([]string{}, false),
	}
}

func (dp *declarationParser) Fetch(input []tokens.Token) (tokens.Token, error) {
	if dp.Index >= len(input) {
		return tokens.Token{}, ParseException{"Reached unexpected end of tokens"}
	}
	tk := input[dp.Index]
	return tk, nil
}

func (dp *declarationParser) ParseConstantState(input []tokens.Token) error {
	descriptor, err := dp.Fetch(input)
	if err != nil {
		return err
	}
	if descriptor.Type != tokens.Identifier {
		return ParseException{"Identifier state descriptor is no keyword identifier"}
	}
	switch descriptor.Value {
	case variableKeyword:
		dp.Declaration.Constant = false
	case constantKeyword:
		dp.Declaration.Constant = true
	default:
		return ParseException{"Unknown identifier state descriptor"}
	}
	dp.Index++
	return nil
}

func (dp *declarationParser) CollectAliases(input []tokens.Token) error {
	for !dp.ExpectAssignment && dp.Index < len(input) {
		active, err := dp.Fetch(input)
		if err != nil {
			return err
		}
		//fmt.Println(active, dp.Index)
		switch active.Type {
		case tokens.Identifier:
			if dp.ExpectTypeInformation {
				for len(dp.Casts) < len(dp.Declaration.Alias) {
					dp.Casts = append(dp.Casts, active.Value)
				}
			} else {
				dp.Declaration.Alias = append(dp.Declaration.Alias, active.Value)
			}
			dp.ExpectTypeInformation = false
		case tokens.Statement:
			return nil
		case tokens.Separator:
		case tokens.Operator:
			if active.Value == ":" {
				dp.ExpectTypeInformation = true
			} else if active.Value == "=" {
				dp.ExpectAssignment = true
			} else {
				return ParseException{"Unexpected operator"}
			}
		default:
			return newUnexpectedTokenException(active.Type)
		}
		dp.Index++
	}
	return nil
}

func (dp *declarationParser) StoreDefaultValues(input []tokens.Token) error {
	if len(dp.Casts) != len(dp.Declaration.Alias) {
		return ParseException{"Required type information not found"}
	}
	for _, t := range dp.Casts {
		dp.Declaration.AddBack(nodes.NewTypecast(t, nodes.NewLiteral(runtime.Value{})))
	}
	return nil
}

func (dp *declarationParser) CollectAssignedValues(input []tokens.Token) error {
	iteration := 0
	for ; dp.Index < len(input); dp.Index++ {
		term, offset, err := newTermParser().Parse(input[dp.Index:])
		if err != nil {
			return err
		}
		dp.Index += offset
		if iteration < len(dp.Casts) {
			term = nodes.NewTypecast(dp.Casts[iteration], term)
			iteration++
		}
		dp.Declaration.AddBack(term)

		// reached end of statement
		if dp.Index < len(input) && input[dp.Index].Type == tokens.Statement {
			return nil
		}
	}
	return nil
}

func (dp *declarationParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	if err := dp.ParseConstantState(input); err != nil {
		return nil, 0, err
	}
	if err := dp.CollectAliases(input); err != nil {
		return nil, 0, err
	}

	// handle if there is no direct assignment
	if !dp.ExpectAssignment {
		if err := dp.StoreDefaultValues(input); err != nil {
			return nil, 0, err
		}
		return dp.Declaration, dp.Index, nil
	}

	// collect values
	if err := dp.CollectAssignedValues(input); err != nil {
		return nil, 0, err
	}
	return dp.Declaration, dp.Index, nil
}
