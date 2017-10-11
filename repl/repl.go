package repl

import (
	"fmt"
	"strings"

	"github.com/tealang/tea-go/lexer"
	"github.com/tealang/tea-go/parser"
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/functions"
	"github.com/tealang/tea-go/runtime/operators"
	"github.com/tealang/tea-go/runtime/types"
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
	fmt.Println("digraph G {\n", strings.Join(ast.Graphviz("head"), "\n"), "}")
	output, err := ast.Eval(r.Context)
	if err != nil {
		return "", err
	}
	if output.Type != nil {
		return output.String(), nil
	}
	return "", nil
}

func (r *Instance) Stop() {
	r.Active = false
}

func New() *Instance {
	ctx := runtime.NewContext()
	operators.Load(ctx)
	types.Load(ctx)
	functions.Load(ctx)
	return &Instance{
		Active:  true,
		Context: ctx,
	}
}
