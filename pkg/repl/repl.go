// Package repl provides a simple interactive runtime environment.
package repl

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/tealang/core/pkg/lexer"
	"github.com/tealang/core/pkg/parser"
	"github.com/tealang/core/pkg/runtime"
	"github.com/tealang/core/pkg/runtime/functions"
	"github.com/tealang/core/pkg/runtime/operators"
	"github.com/tealang/core/pkg/runtime/types"
)

// Config stores the REPL instance configuration.
type Config struct {
	OutputGraph bool
}

// Instance is a REPL runtime instance.
type Instance struct {
	Active  bool
	context *runtime.Context
	cfg     Config
}

const (
	graphvizFormat = "digraph G {\n%s\n}"
	graphvizItem   = "head"
)

// Load fetches a file from disk and evaluates it in the runtime instance.
func (r *Instance) Load(file string) error {
	code, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "can not execute")
	}
	_, err = r.Interpret(string(code))
	if err != nil {
		return errors.Wrap(err, "execution failed")
	}
	return nil
}

// Interpret runs the given input program in the runtime instance.
func (r *Instance) Interpret(input string) (string, error) {
	tokens := lexer.Lex(input)
	ast, _, err := parser.Parse(tokens)
	if err != nil {
		return "", errors.Wrap(err, "failed to interpret")
	}
	if r.cfg.OutputGraph {
		return fmt.Sprintf(graphvizFormat, strings.Join(ast.Graphviz(graphvizItem), "\n")), nil
	}
	output, err := ast.Eval(r.context)
	if err != nil {
		return "", errors.Wrap(err, "failed to interpret")
	}
	if output.Type != nil {
		return output.String(), nil
	}
	return "", nil
}

// Stop kills the running instance.
func (r *Instance) Stop() {
	r.Active = false
}

// New instantiates a new REPL runtime instance.
func New(cfg Config) *Instance {
	ctx := runtime.NewContext()

	// load language runtime
	operators.Load(ctx)
	types.Load(ctx)
	functions.Load(ctx)

	return &Instance{
		Active:  true,
		context: ctx,
		cfg:     cfg,
	}
}
