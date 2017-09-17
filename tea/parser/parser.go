package parser

import (
	"errors"
	"fmt"

	"github.com/tealang/tea-go/tea/lexer/tokens"
	"github.com/tealang/tea-go/tea/runtime"
	"github.com/tealang/tea-go/tea/runtime/nodes"
)

type SequenceLevel int

const (
	LevelGlobal SequenceLevel = iota
	LevelFunction
)

type ParseException struct {
	Message string
}

func (p ParseException) Error() string {
	return fmt.Sprintf("ParseException: %s", p.Message)
}

// Parse generates an abstract syntax tree from the given list of tokens.
func Parse(input []tokens.Token) (nodes.Node, error) {
	seq, _, err := GenerateSequence(input, false, LevelGlobal)
	if err != nil {
		return nil, err
	}
	return seq, nil
}

func GenerateTerm(input []tokens.Token) (nodes.Node, int, error) {
	return nil, 0, nil
}

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

func GenerateReturn(input []tokens.Token) (nodes.Node, int, error) {
	ctrl := nodes.NewController(runtime.BehaviorReturn)
	if len(input) < 2 {
		return ctrl, 1, nil
	}

	term, n, err := GenerateTerm(input[1:])
	if err != nil {
		return ctrl, n + 1, err
	}
	ctrl.AddBack(term)
	return ctrl, n + 1, nil
}

func GenerateSequence(input []tokens.Token, block bool, level SequenceLevel) (nodes.Node, int, error) {
	var (
		index = 0
		seq   = nodes.NewSequence(block)
	)
	for ; index < len(input); index++ {
		requireEndStatement := true
		active := input[index]
		switch active.Type {
		case tokens.Identifier:
			switch active.Value {
			case "let", "var":
				stmt, n, err := GenerateDeclaration(input[index:])
				if err != nil {
					return seq, index, err
				}
				seq.AddBack(stmt)
				index += n
			case "return":
				stmt, n, err := GenerateReturn(input[index:])
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
			}
		}

		if index >= len(input) {
			requireEndStatement = false
		}

		if requireEndStatement && input[index].Type != tokens.Statement {
			return seq, index, errors.New("Statement end required")
		}
	}
	return seq, index, nil
}

func GetOperatorPriority(symbol string, previous tokens.Token) (int, error) {
	switch symbol {
	case "&", "|":
		return 8, nil
	case "!":
		return 7, nil
	case "^":
		return 6, nil
	case "*", "/":
		return 5, nil
	case "+", "-", ":":
		return 4, nil
	case "%":
		return 3, nil
	case "<", ">", ">=", "<=", "=>", "=<", "!=", "==":
		return 2, nil
	case "&&", "||", "^|":
		return 1, nil
	}
	return 0, nil
}

func CleanWhitespace(input []tokens.Token) []tokens.Token {
	output := make([]tokens.Token, 0, len(input))
	for _, tk := range input {
		if tk.Type != tokens.Whitespace {
			output = append(output, tk)
		}
	}
	return output
}
