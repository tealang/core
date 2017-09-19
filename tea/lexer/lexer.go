package lexer

import "github.com/tealang/tea-go/tea/lexer/tokens"

// Lex converts the input into a series of tokens.
func Lex(input string) []tokens.Token {
	active := tokens.Token{
		Value: "",
		Type:  nil,
	}
	output := []tokens.Token{}
	for _, c := range input {
		value := active.Value + string(c)
		if active.Type != nil && active.Type.Match(value) {
			active.Value = value
		} else {
			if active.Type != nil {
				output = append(output, active)
			}
			active = tokens.Token{
				Value: string(c),
				Type:  tokens.FindMatch(string(c)),
			}
		}
	}
	return append(output, active)
}
