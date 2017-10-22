package parser

import (
	"fmt"

	"github.com/tealang/core/lexer/tokens"
)

// Exception is a exception that occured while parsing the input.
type Exception struct {
	message string
}

func (p Exception) Error() string {
	return fmt.Sprintf("error while parsing: %s", p.message)
}

func newParseException(message string) Exception {
	return Exception{message}
}

func newUnexpectedTokenException(t tokens.Token) Exception {
	return newParseException(fmt.Sprintf("did not expect token %s of type %s", t.Value, t.Type.Name))
}

func newMissingOperatorException() Exception {
	return newParseException(fmt.Sprintf("missing operator"))
}

func newMissingOperandsException(expected, got int) Exception {
	return newParseException(fmt.Sprintf("not enough operands, expected %d, got %d", expected, got))
}
