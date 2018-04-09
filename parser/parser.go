package parser

import (
	"github.com/pkg/errors"
	"github.com/tealang/core/lexer/tokens"
	"github.com/tealang/core/runtime/nodes"
)

const (
	constantKeyword    = "let"
	variableKeyword    = "var"
	returnKeyword      = "return"
	breakKeyword       = "break"
	fallthroughKeyword = "fallthrough"
	defaultKeyword     = "default"
	continueKeyword    = "continue"
	ifKeyword          = "if"
	elseKeyword        = "else"
	forKeyword         = "for"
	functionKeyword    = "func"
	operatorKeyword    = "operator"
	trueKeyword        = "true"
	falseKeyword       = "false"
	nullKeyword        = "null"
	matchKeyword       = "match"
	caseKeyword        = "case"
	castOperator       = ":"
	assignmentOperator = "="
)

// Parse generates an abstract syntax tree from the given list of tokens.
// It returns the generated tree node, the parsed token offset and in the case of a failure,
// an error object.
func Parse(input []tokens.Token) (nodes.Node, int, error) {
	// clean input from whitespace
	cleaned := make([]tokens.Token, 0, len(input))
	for _, tk := range input {
		if tk.Type != tokens.Whitespace {
			cleaned = append(cleaned, tk)
		}
	}

	seq, n, err := newSequenceParser(false, 0).Parse(cleaned)
	if err != nil {
		return nil, n, errors.Wrap(err, "error while parsing")
	}
	return seq, n, nil
}
