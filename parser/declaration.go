package parser

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/lexer/tokens"
	"github.com/tealang/core/runtime"
	"github.com/tealang/core/runtime/nodes"
)

type declarationParser struct {
	assignment  bool
	typecasts   []string
	declaration *nodes.Declaration
	index, size int
	input       []tokens.Token
}

func newDeclarationParser() *declarationParser {
	return &declarationParser{
		typecasts:   make([]string, 0),
		declaration: nodes.NewMultiDeclaration([]string{}, false),
	}
}

func (dp *declarationParser) fetch() (tokens.Token, error) {
	if dp.index >= dp.size {
		return tokens.Token{}, errors.New("unexpected end of tokens while fetching")
	}
	tk := dp.input[dp.index]
	return tk, nil
}

func (dp *declarationParser) parseMode() error {
	descriptor, err := dp.fetch()
	if err != nil {
		return err
	}
	if descriptor.Type != tokens.Identifier {
		return errors.Errorf("expected state descriptor, got %s", descriptor.Type)
	}
	switch descriptor.Value {
	case variableKeyword:
		dp.declaration.Constant = false
	case constantKeyword:
		dp.declaration.Constant = true
	default:
		return errors.New("state descriptor must be either let or var")
	}
	dp.index++
	return nil
}

func (dp *declarationParser) collectAliases() error {
	nextTypeInfo := false
	for !dp.assignment && dp.index < dp.size {
		active, err := dp.fetch()
		if err != nil {
			return err
		}
		//fmt.Println(active, dp.Index)
		switch active.Type {
		case tokens.Identifier:
			if nextTypeInfo {
				for len(dp.typecasts) < len(dp.declaration.Alias) {
					dp.typecasts = append(dp.typecasts, active.Value)
				}
			} else {
				dp.declaration.Alias = append(dp.declaration.Alias, active.Value)
			}
			nextTypeInfo = false
		case tokens.Statement:
			return nil
		case tokens.Separator:
		case tokens.Operator:
			switch active.Value {
			case castOperator:
				nextTypeInfo = true
			case assignmentOperator:
				dp.assignment = true
			default:
				return errors.Errorf("did expect typecast or assignment operator, got %s", active.Value)
			}
		default:
			return errors.Errorf("did not expect token %s", active.Type)
		}
		dp.index++
	}
	return nil
}

func (dp *declarationParser) assignDefaultValues() error {
	if len(dp.typecasts) != len(dp.declaration.Alias) {
		return errors.Errorf("expected %d typecasts, got %d", len(dp.declaration.Alias), len(dp.typecasts))
	}
	for _, t := range dp.typecasts {
		dp.declaration.AddBack(nodes.NewTypecast(t, nodes.NewLiteral(runtime.Value{})))
	}
	return nil
}

func (dp *declarationParser) assignValues() error {
	for i := 0; i < len(dp.declaration.Alias); i++ {
		term, n, err := newTermParser().Parse(dp.input[dp.index:])
		if err != nil {
			return err
		}
		dp.index += n

		if i < len(dp.typecasts) {
			term = nodes.NewTypecast(dp.typecasts[i], term)
		}
		dp.declaration.AddBack(term)

		if dp.index < dp.size && dp.input[dp.index].Type == tokens.Separator {
			dp.index++
		}
	}
	return nil
}

func (dp *declarationParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	dp.index, dp.size = 0, len(input)
	dp.input = input

	if err := dp.parseMode(); err != nil {
		return nil, 0, errors.Wrap(err, "failed to parse state")
	}
	if err := dp.collectAliases(); err != nil {
		return nil, 0, errors.Wrap(err, "failed collecting alias")
	}

	// handle if there is no direct assignment
	if !dp.assignment {
		if err := dp.assignDefaultValues(); err != nil {
			return nil, 0, errors.Wrap(err, "could not assign null values")
		}
		return dp.declaration, dp.index, nil
	}

	// collect values
	if err := dp.assignValues(); err != nil {
		return nil, 0, errors.Wrap(err, "failed collecting values")
	}
	return dp.declaration, dp.index, nil
}
