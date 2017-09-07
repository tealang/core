package runtime

type ContextBehavior int

const (
	BehaviorDefault ContextBehavior = iota
)

type Context struct {
	Namespace       *Namespace
	GlobalNamespace *Namespace
	Behavior        ContextBehavior
}

func (c *Context) Substitute(f func(*Context) (Value, error)) (Value, error) {
	backup := c.Namespace
	defer func() { c.Namespace = backup }()
	return f(c)
}
