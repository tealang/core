package tokens

import (
	"fmt"
	"regexp"
)

type TokenMatcher func(s string) bool

func NewTokenMatcher(expr string) TokenMatcher {
	regex, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	return func(s string) bool {
		return regex.Match([]byte(s))
	}
}

type TokenType struct {
	Name  string
	Match TokenMatcher
}

func (tt TokenType) String() string {
	return tt.Name
}

type Token struct {
	Type  *TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("%s('%s')", t.Type, t.Value)
}

var (
	LeftParentheses = &TokenType{
		Name:  "leftParentheses",
		Match: NewTokenMatcher(`^\($`),
	}
	RightParentheses = &TokenType{
		Name:  "rightParentheses",
		Match: NewTokenMatcher(`^\)$`),
	}
	Operator = &TokenType{
		Name:  "operator",
		Match: NewTokenMatcher(`^([+\-*/=:<>!%^&|]|([+\-*/^%<>=!]=)|([|^]\|)|(&&))$`),
	}
	Whitespace = &TokenType{
		Name:  "whitespace",
		Match: NewTokenMatcher(`^\s+$`),
	}
	Number = &TokenType{
		Name:  "number",
		Match: NewTokenMatcher(`^\-?[0-9]+(\.[0-9]*)?$`),
	}
	Identifier = &TokenType{
		Name:  "identifier",
		Match: NewTokenMatcher(`^(#|[a-zA-Z_])+([0-9a-zA-Z_]+)?$`),
	}
	String = &TokenType{
		Name:  "string",
		Match: NewTokenMatcher(`^"(\\(["abfnrtv])?|[^\n\r"])*"?$`),
	}
	Statement = &TokenType{
		Name:  "statement",
		Match: NewTokenMatcher(`^;$`),
	}
	Separator = &TokenType{
		Name:  "separator",
		Match: NewTokenMatcher(`^,$`),
	}
	RightBlock = &TokenType{
		Name:  "rightBlock",
		Match: NewTokenMatcher(`^}$`),
	}
	LeftBlock = &TokenType{
		Name:  "leftBlock",
		Match: NewTokenMatcher(`^{$`),
	}
	AllTypes = []*TokenType{
		LeftParentheses,
		RightParentheses,
		Operator,
		Whitespace,
		Separator,
		Number,
		Identifier,
		String,
		Statement,
		RightBlock,
		LeftBlock,
	}
)

func FindMatch(value string) *TokenType {
	for _, tt := range AllTypes {
		if tt.Match(value) {
			return tt
		}
	}
	return nil
}
