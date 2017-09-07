// Package tea provides an implementation of the Tealang runtime.
package tea

type Evaluable interface {
	Eval(c *Context) (Value, error)
}
