// Package runtime provides the Tealang runtime functionality including namespaces, context and execution nodes.
package runtime

// ContextBehavior controls the runtime execution behavior.
type ContextBehavior int

const (
	// BehaviorDefault is the normal behavior.
	BehaviorDefault ContextBehavior = iota
	// BehaviorContinue skips the body of the loop.
	BehaviorContinue
	// BehaviorBreak breaks the loop.
	BehaviorBreak
	// BehaviorReturn cancels the running sequence up to function level.
	BehaviorReturn
	// BehaviorFallthrough executes the next case in the match list.
	BehaviorFallthrough
)

// Context is the runtime context the AST is executed in.
type Context struct {
	Namespace       *Namespace
	GlobalNamespace *Namespace
	Behavior        ContextBehavior
}

// Substitute executes the method in a substituted namespace.
// This makes new variables within the substituted namespace unavailable to everyone outside of it.
func (c *Context) Substitute(f func(*Context) (Value, error)) (Value, error) {
	backup := c.Namespace
	c.Namespace = NewNamespace(c.Namespace)
	defer func() { c.Namespace = backup }()
	return f(c)
}

// NewContext instantiates a new context with only one namespace, both local and global.
func NewContext() *Context {
	ns := NewNamespace(nil)
	return &Context{
		Namespace:       ns,
		GlobalNamespace: ns,
		Behavior:        BehaviorDefault,
	}
}
