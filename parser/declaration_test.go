package parser

import (
	"testing"

	"github.com/tealang/core/lexer/tokens"
)

func Test_declarationParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   []tokens.Token
		want    int
		wantErr bool
	}{
		{
			"Single declaration",
			[]tokens.Token{
				{Type: tokens.Identifier, Value: "let"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Operator, Value: ":"},
				{Type: tokens.Identifier, Value: "int"},
				{Type: tokens.Statement, Value: ";"},
			},
			4,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := newDeclarationParser()
			_, n, err := dp.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("declarationParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if n != tt.want {
				t.Errorf("declarationParser.Parse() n = %v, want %v", n, tt.want)
			}
		})
	}
}
