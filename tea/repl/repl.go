package repl

import (
	"fmt"

	"github.com/tealang/tea-go/tea/lexer"
)

type Instance struct {
	Active bool
}

func (r *Instance) Interpret(input string) string {
	return fmt.Sprint(lexer.Lex(input))
}

func (r *Instance) Stop() {
	r.Active = false
}

func New() *Instance {
	return &Instance{
		Active: true,
	}
}
