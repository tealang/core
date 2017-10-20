package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime/nodes"
)

const (
	constantKeyword = "let"
	variableKeyword = "var"
	returnKeyword   = "return"
	breakKeyword    = "break"
	defaultKeyword  = "default"
	continueKeyword = "continue"
	ifKeyword       = "if"
	elseKeyword     = "else"
	forKeyword      = "for"
	functionKeyword = "func"
	trueKeyword     = "true"
	falseKeyword    = "false"
	nullKeyword     = "null"
)

// Parser provides an interface for token parsing.
type Parser interface {
	// Parse generates an abstract syntax tree from the given list of tokens.
	// It returns the generated tree node, the parsed token offset and in the case of a failure,
	// an error object.
	Parse(input []tokens.Token) (nodes.Node, int, error)
}

// New instantiates a new program parser.
func New() *programParser {
	return &programParser{}
}

type programParser struct {
}

func (pp programParser) Parse(input []tokens.Token) (nodes.Node, error) {
	seq, _, err := newSequenceParser(false, 0).Parse(pp.clean(input))
	if err != nil {
		return nil, err
	}
	return seq, nil
}

func (programParser) clean(input []tokens.Token) []tokens.Token {
	output := make([]tokens.Token, 0, len(input))
	for _, tk := range input {
		if tk.Type != tokens.Whitespace {
			output = append(output, tk)
		}
	}
	return output
}
