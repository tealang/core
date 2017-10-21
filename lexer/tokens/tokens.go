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

type Type struct {
	Name  string
	Match TokenMatcher
}

func (tt Type) String() string {
	return tt.Name
}

type Token struct {
	Type  *Type
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("%s('%s')", t.Type, t.Value)
}

var (
	LeftParentheses = &Type{
		Name:  "leftParentheses",
		Match: NewTokenMatcher(`^\($`),
	}
	RightParentheses = &Type{
		Name:  "rightParentheses",
		Match: NewTokenMatcher(`^\)$`),
	}
	Operator = &Type{
		Name:  "operator",
		Match: NewTokenMatcher(`^([+\-*/=:<>!%^&|]|([+\-*/^%<>=!][=?]{1})|([|^]\|)|(&&))$`),
	}
	Whitespace = &Type{
		Name:  "whitespace",
		Match: NewTokenMatcher(`^\s+$`),
	}
	Number = &Type{
		Name:  "number",
		Match: NewTokenMatcher(`^\-?[0-9]+(\.[0-9]*)?$`),
	}
	Identifier = &Type{
		Name:  "identifier",
		Match: NewTokenMatcher(`^(#|[a-zA-Z_])+([0-9a-zA-Z_]+)?$`),
	}
	String = &Type{
		Name:  "string",
		Match: NewTokenMatcher(`^"(\\(["abfnrtv])?|[^\n\r"])*"?$`),
	}
	Statement = &Type{
		Name:  "statement",
		Match: NewTokenMatcher(`^;$`),
	}
	Separator = &Type{
		Name:  "separator",
		Match: NewTokenMatcher(`^,$`),
	}
	RightBlock = &Type{
		Name:  "rightBlock",
		Match: NewTokenMatcher(`^}$`),
	}
	LeftBlock = &Type{
		Name:  "leftBlock",
		Match: NewTokenMatcher(`^{$`),
	}
	AllTypes = []*Type{
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

func FindMatch(value string) *Type {
	for _, tt := range AllTypes {
		if tt.Match(value) {
			return tt
		}
	}
	return nil
}
