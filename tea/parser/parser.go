package parser

import (
	"fmt"

	"github.com/tealang/tea-go/tea/lexer/tokens"
	"github.com/tealang/tea-go/tea/runtime/nodes"
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
