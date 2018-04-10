// Package lexer provides a simple lexer that generates a discrete series of tokens from an input string.
package lexer

import "github.com/tealang/core/pkg/lexer/tokens"

// Lex converts the input into a series of tokens.
func Lex(input string) []tokens.Token {
	active := tokens.Token{
		Value: "",
		Type:  nil,
	}
	output := []tokens.Token{}
	for i := 0; i < len(input); i++ {
		c := input[i]
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
			switch active.Type {
			case tokens.SingleLineComment:
				for i < len(input) && input[i] != '\n' {
					i++
				}
				active = tokens.Token{Type: tokens.Whitespace}
			}
		}
	}
	return append(output, active)
}
