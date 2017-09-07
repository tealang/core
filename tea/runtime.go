package tea

import (
	"fmt"
)

type Runtime struct {
	Active bool
}

func (r *Runtime) Interpret(input string) string {
	return fmt.Sprint(Lex(input))
}

func (r *Runtime) Stop() {
	r.Active = false
}

func NewRuntime() *Runtime {
	return &Runtime{
		Active: true,
	}
}
