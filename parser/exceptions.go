package parser

import (
	"fmt"

	"github.com/tealang/tea-go/lexer/tokens"
)

type ParseException struct {
	message string
}

func (p ParseException) Error() string {
	return fmt.Sprintf("ParseException: %s", p.message)
}

func newParseException(message string) ParseException {
	return ParseException{message}
}

func newUnexpectedTokenException(t *tokens.Type) ParseException {
	return newParseException(fmt.Sprintf("Did not expect token %s", t.Name))
}

func newMissingOperatorException() ParseException {
	return newParseException(fmt.Sprintf("Missing operator"))
}

func newMissingOperandsException(expected, got int) ParseException {
	return newParseException(fmt.Sprintf("Missing operands, expected %d, got %d", expected, got))
}
