package runtime

import (
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()
	if ctx.GlobalNamespace != ctx.Namespace {
		t.Error("Global and local namespace should be equal")
	}
	if ctx.Behavior != BehaviorDefault {
		t.Error("Behavior on New should be default")
	}
}

func TestContext_Substitute(t *testing.T) {
	ctx := NewContext()
	v := Value{Name: "v"}
	if x, err := ctx.Substitute(func(c *Context) (Value, error) {
		c.Namespace.Store(v)
		return v, nil
	}); err != nil {
		t.Errorf("Unexpected error on substitute: %s", err)
	} else if !reflect.DeepEqual(x, v) {
		t.Errorf("Invalid value: expected %v, got %v", v, x)
	}
	if _, err := ctx.Namespace.Find(SearchIdentifier, "v"); err == nil {
		t.Errorf("Expected error on find: %s", err)
	}
	if _, err := ctx.GlobalNamespace.Find(SearchIdentifier, "v"); err == nil {
		t.Errorf("Expected error on find (global): %s", err)
	}
}
