package parser

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/lexer/tokens"
	"github.com/tealang/core/runtime"
	"github.com/tealang/core/runtime/nodes"
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
		return tokens.Token{}, errors.New("unexpected end of tokens while fetching")
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
		return errors.Errorf("expected state descriptor, got %s", descriptor.Type)
	}
	switch descriptor.Value {
	case variableKeyword:
		dp.Declaration.Constant = false
	case constantKeyword:
		dp.Declaration.Constant = true
	default:
		return errors.New("state descriptor must be either let or var")
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
				return errors.Errorf("did expect typecast or assignment operator, got %s", active.Value)
			}
		default:
			return errors.Errorf("did not expect token %s", active.Type)
		}
		dp.Index++
	}
	return nil
}

func (dp *declarationParser) StoreDefaultValues(input []tokens.Token) error {
	if len(dp.Casts) != len(dp.Declaration.Alias) {
		return errors.Errorf("expected %d typecasts, got %d", len(dp.Declaration.Alias), len(dp.Casts))
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
		return nil, 0, errors.Wrap(err, "failed to parse state")
	}
	if err := dp.CollectAliases(input); err != nil {
		return nil, 0, errors.Wrap(err, "failed collecting alias")
	}

	// handle if there is no direct assignment
	if !dp.ExpectAssignment {
		if err := dp.StoreDefaultValues(input); err != nil {
			return nil, 0, errors.Wrap(err, "could not assign null values")
		}
		return dp.Declaration, dp.Index, nil
	}

	// collect values
	if err := dp.CollectAssignedValues(input); err != nil {
		return nil, 0, errors.Wrap(err, "failed collecting values")
	}
	return dp.Declaration, dp.Index, nil
}
