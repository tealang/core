package lexer

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
	AllTypes = []TokenType{
		{
			Name:  "leftParentheses",
			Match: NewTokenMatcher(`^\($`),
		},
		{
			Name:  "rightParentheses",
			Match: NewTokenMatcher(`^\)$`),
		},
		{
			Name:  "operator",
			Match: NewTokenMatcher(`^([+\-*/=:<>!%^&|]|([+\-*/^%<>=!]=)|([|^]\|)|(&&))$`),
		},
		{
			Name:  "whitespace",
			Match: NewTokenMatcher(`^\s+$`),
		},
		{
			Name:  "number",
			Match: NewTokenMatcher(`^\-?[0-9]+(\.[0-9]*)?$`),
		},
		{
			Name:  "identifier",
			Match: NewTokenMatcher(`^(#|[a-zA-Z_])+([0-9a-zA-Z_]+)?$`),
		},
		{
			Name:  "string",
			Match: NewTokenMatcher(`^"(\\(["abfnrtv])?|[^\n\r"])*"?$`),
		},
		{
			Name:  "statement",
			Match: NewTokenMatcher(`^;$`),
		},
		{
			Name:  "separator",
			Match: NewTokenMatcher(`^,$`),
		},
		{
			Name:  "rightBlock",
			Match: NewTokenMatcher(`^}$`),
		},
		{
			Name:  "leftBlock",
			Match: NewTokenMatcher(`^{$`),
		},
	}
)

func FindMatch(value string) *TokenType {
	for _, tt := range AllTypes {
		if tt.Match(value) {
			return &tt
		}
	}
	return nil
}

func Lex(input string) []Token {
	active := Token{
		Value: "",
		Type:  nil,
	}
	output := []Token{}
	for _, c := range input {
		value := active.Value + string(c)
		if active.Type != nil && active.Type.Match(value) {
			active.Value = value
		} else {
			if active.Type != nil {
				output = append(output, active)
			}
			active = Token{
				Value: string(c),
				Type:  FindMatch(string(c)),
			}
		}
	}
	return append(output, active)
}
