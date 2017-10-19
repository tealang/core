package repl

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/tealang/tea-go/lexer"
	"github.com/tealang/tea-go/parser"
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/functions"
	"github.com/tealang/tea-go/runtime/operators"
	"github.com/tealang/tea-go/runtime/types"
)

type Config struct {
	OutputGraph bool
}

type Instance struct {
	context *runtime.Context
	active  bool
	cfg     Config
}

const (
	graphvizFormat = "digraph G {\n%s\n}"
	graphvizItem   = "head"
)

func (r *Instance) IsActive() bool {
	return r.active
}

func (r *Instance) Load(file string) error {
	code, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	_, err = r.Interpret(string(code))
	if err != nil {
		return err
	}
	return nil
}

func (r *Instance) Interpret(input string) (string, error) {
	tokens := lexer.Lex(input)
	ast, err := parser.New().Parse(tokens)
	if err != nil {
		return "", err
	}
	if r.cfg.OutputGraph {
		return fmt.Sprintf(graphvizFormat, strings.Join(ast.Graphviz(graphvizItem), "\n")), nil
	}
	output, err := ast.Eval(r.context)
	if err != nil {
		return "", err
	}
	if output.Type != nil {
		return output.String(), nil
	}
	return "", nil
}

func (r *Instance) Stop() {
	r.active = false
}

func New(cfg Config) *Instance {
	ctx := runtime.NewContext()

	// load language runtime
	operators.Load(ctx)
	types.Load(ctx)
	functions.Load(ctx)

	return &Instance{
		active:  true,
		context: ctx,
		cfg:     cfg,
	}
}
