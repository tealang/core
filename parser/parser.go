package parser

import (
	"fmt"

	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime/nodes"
)

const (
	KeywordConstant = "let"
	KeywordVariable = "var"
)

type Parser interface {
	// Parse generates an abstract syntax tree from the given list of tokens.
	// It returns the generated tree node, the parsed token offset and in the case of a failure,
	// an error object.
	Parse(input []tokens.Token) (nodes.Node, int, error)
}

type ParseException struct {
	Message string
}

func (p ParseException) Error() string {
	return fmt.Sprintf("ParseException: %s", p.Message)
}

type ProgramParser struct {
}

func (ProgramParser) Parse(input []tokens.Token) (nodes.Node, error) {
	seq, _, err := NewSequenceParser().Parse(input, false, LevelGlobal)
	if err != nil {
		return nil, err
	}
	return seq, nil
}

func (ProgramParser) CleanWhitespace(input []tokens.Token) []tokens.Token {
	output := make([]tokens.Token, 0, len(input))
	for _, tk := range input {
		if tk.Type != tokens.Whitespace {
			output = append(output, tk)
		}
	}
	return output
}
