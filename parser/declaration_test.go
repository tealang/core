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
			"Single declaration, invalid mode",
			[]tokens.Token{
				{Type: tokens.Identifier, Value: "invalid"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Operator, Value: ":"},
				{Type: tokens.Identifier, Value: "int"},
				{Type: tokens.Statement, Value: ";"},
			},
			0,
			true,
		},
		{
			"Single declaration, invalid mode token",
			[]tokens.Token{
				{Type: tokens.Number, Value: "3412"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Operator, Value: ":"},
				{Type: tokens.Identifier, Value: "int"},
				{Type: tokens.Statement, Value: ";"},
			},
			0,
			true,
		},
		{
			"Unexpected end of tokens",
			[]tokens.Token{

				{Type: tokens.Identifier, Value: "invalid"},
				{Type: tokens.Identifier, Value: "x"},
			},
			0,
			true,
		},
		{
			"Unexpected tokens",
			[]tokens.Token{

				{Type: tokens.Identifier, Value: "let"},
				{Type: tokens.Number, Value: "3"},
			},
			0,
			true,
		},
		{
			"Unexpected operator",
			[]tokens.Token{

				{Type: tokens.Identifier, Value: "let"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Operator, Value: "?"},
				{Type: tokens.Identifier, Value: "int"},
				{Type: tokens.Statement, Value: ";"},
			},
			0,
			true,
		},
		{
			"Single declaration (explicit w/o assignment)",
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
		{
			"Tuple declaration (explicit w/o assignment)",
			[]tokens.Token{
				{Type: tokens.Identifier, Value: "let"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Separator},
				{Type: tokens.Identifier, Value: "y"},
				{Type: tokens.Operator, Value: ":"},
				{Type: tokens.Identifier, Value: "int"},
				{Type: tokens.Statement, Value: ";"},
			},
			6,
			false,
		},
		{
			"Tuple declaration (explicit w/o assignment, different types)",
			[]tokens.Token{
				{Type: tokens.Identifier, Value: "let"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Operator, Value: ":"},
				{Type: tokens.Identifier, Value: "float"},
				{Type: tokens.Separator},
				{Type: tokens.Identifier, Value: "y"},
				{Type: tokens.Operator, Value: ":"},
				{Type: tokens.Identifier, Value: "int"},
				{Type: tokens.Statement, Value: ";"},
			},
			8,
			false,
		},
		{
			"Single declaration (explicit with assignment)",
			[]tokens.Token{
				{Type: tokens.Identifier, Value: "var"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Operator, Value: ":"},
				{Type: tokens.Identifier, Value: "int"},
				{Type: tokens.Operator, Value: "="},
				{Type: tokens.Number, Value: "3"},
				{Type: tokens.Statement, Value: ";"},
			},
			6,
			false,
		},
		{
			"Tuple declaration (explicit with assignment)",
			[]tokens.Token{
				{Type: tokens.Identifier, Value: "var"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Separator},
				{Type: tokens.Identifier, Value: "y"},
				{Type: tokens.Operator, Value: ":"},
				{Type: tokens.Identifier, Value: "int"},
				{Type: tokens.Operator, Value: "="},
				{Type: tokens.Number, Value: "3"},
				{Type: tokens.Separator},
				{Type: tokens.Number, Value: "4"},
				{Type: tokens.Statement, Value: ";"},
			},
			10,
			false,
		},
		{
			"Single declaration (implicit with assignment)",
			[]tokens.Token{
				{Type: tokens.Identifier, Value: "let"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Operator, Value: "="},
				{Type: tokens.Number, Value: "3"},
				{Type: tokens.Statement, Value: ";"},
			},
			4,
			false,
		},
		{
			"Tuple declaration (implicit with assignment)",
			[]tokens.Token{
				{Type: tokens.Identifier, Value: "let"},
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Separator},
				{Type: tokens.Identifier, Value: "y"},
				{Type: tokens.Operator, Value: "="},
				{Type: tokens.Number, Value: "3"},
				{Type: tokens.Separator},
				{Type: tokens.Number, Value: "4"},
				{Type: tokens.Statement, Value: ";"},
			},
			8,
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
