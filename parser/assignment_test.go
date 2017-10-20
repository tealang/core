package parser

import (
	"testing"

	"github.com/tealang/core/lexer/tokens"
)

func Test_assignmentParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   []tokens.Token
		want    int
		wantErr bool
	}{
		{
			"Tuple assignment",
			[]tokens.Token{
				{Type: tokens.Identifier, Value: "x"},
				{Type: tokens.Separator},
				{Type: tokens.Identifier, Value: "y"},
				{Type: tokens.Operator, Value: "="},
				{Type: tokens.Number, Value: "1"},
				{Type: tokens.Number, Value: "2"},
				{Type: tokens.Statement},
			},
			6,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := newAssignmentParser()
			_, n, err := ap.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("assignmentParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if n != tt.want {
				t.Errorf("assignmentParser.Parse() n = %v, want %v", n, tt.want)
			}
		})
	}
}
