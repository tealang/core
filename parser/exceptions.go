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

func newUnexpectedTokenException(t *tokens.TokenType) ParseException {
	return newParseException(fmt.Sprintf("UnexpectedTokenException: Did not expect token %s", t.Name))
}