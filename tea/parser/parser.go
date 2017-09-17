package parser

import (
	"github.com/tealang/tea-go/tea/lexer/tokens"
	"github.com/tealang/tea-go/tea/runtime/nodes"
)

type SequenceLevel int

const (
	LevelGlobal SequenceLevel = iota
	LevelFunction
)

// Parse generates an abstract syntax tree from the given list of tokens.
func Parse(input []tokens.Token) (nodes.Node, error) {
	return nodes.NewSequence(false), nil
}

func GenerateSequence(input []tokens.Token, level SequenceLevel) (nodes.Node, error) {
	return nodes.NewSequence(false), nil
}

func CleanWhitespace(input []tokens.Token) []tokens.Token {
	output := make([]tokens.Token, 0, len(input))
	for _, tk := range input {
		if tk.Type != &tokens.Whitespace {
			output = append(output, tk)
		}
	}
	return output
}
