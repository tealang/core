package runtime

type Evaluable interface {
	Eval(c *Context) (Value, error)
}
