package repl

import (
	"github.com/tealang/tea-go/lexer"
	"github.com/tealang/tea-go/parser"
	"github.com/tealang/tea-go/runtime"
)

type Instance struct {
	Context *runtime.Context
	Active  bool
}

func (r *Instance) Interpret(input string) (string, error) {
	tokens := lexer.Lex(input)
	ast, err := parser.New().Parse(tokens)
	if err != nil {
		return "", err
	}
	output, err := ast.Eval(r.Context)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

func (r *Instance) Stop() {
	r.Active = false
}

func New() *Instance {
	return &Instance{
		Active:  true,
		Context: runtime.NewContext(),
	}
}
