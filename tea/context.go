package tea

import (
	"fmt"
)

type Signature struct {
	Expected []Value
	Function Evaluable
}

func (sign *Signature) Match(args []Value) ([]Value, error) {
	expected, got := len(sign.Expected), len(args)
	if expected < got {
		return nil, ArgumentException(expected, got)
	}

	matched := make([]Value, expected)
	for i := range sign.Expected {
		if got > i {
			if !args[i].Type.KindOf(sign.Expected[i].Type) {
				return nil, ArgumentCastException(sign.Expected[i].Type, args[i].Type)
			}
			casted, err := sign.Expected[i].Type.Cast(args[i])
			if err != nil {
				return nil, err
			}
			matched[i] = casted
		} else if sign.Expected[i].Data != nil {
			matched[i], _ = sign.Expected[i].Type.Cast(sign.Expected[i])
		} else {
			return nil, ArgumentException(expected, got)
		}
		matched[i].Name = sign.Expected[i].Name
	}

	return matched, nil
}

func (sign *Signature) String() string {
	return fmt.Sprintf("<S {%s}>", sign.Expected)
}

type Function struct {
	Signatures []Signature
	Name       string
	Source     *Namespace
}

func (f *Function) Eval(args []Value, c *Context) (Value, error) {
	for _, sign := range f.Signatures {
		matched, err := sign.Match(args)
		if err != nil {
			continue
		}
		return c.Substitute(func(c *Context) (Value, error) {
			executionSpace := NewNamespace(f.Source)
			for _, arg := range matched {
				executionSpace.Store(arg)
			}
			return sign.Function.Eval(c)
		})
	}
	return Value{}, FunctionException(f.Name)
}

type Formatter func(v Value) string
type Caster func(v Value) (Value, error)

type Datatype struct {
	Name   string
	Parent *Datatype
	Cast   Caster
	Format Formatter
}

func (datatype *Datatype) KindOf(other *Datatype) bool {
	if datatype != other {
		return false
	} else if datatype.Parent != nil {
		return datatype.Parent.KindOf(other)
	}
	return true
}

func (datatype *Datatype) String() string {
	return fmt.Sprintf("<T %s{%s}>", datatype.Name, datatype.Parent)
}

type Value struct {
	Type     *Datatype
	Data     interface{}
	Name     string
	Constant bool
}

func (v Value) Alias() string {
	return v.Name
}

func (v Value) SearchSpace() SearchSpace {
	return IdentifierSearchSpace
}

func (v Value) Update(item SearchItem) (SearchItem, error) {
	if v.SearchSpace() != item.SearchSpace() {
		return v, SearchSpaceException()
	}
	switch c := item.(type) {
	case Value:
		if v.Constant {
			return v, ConstantException(v.Name)
		}
		if !c.Type.KindOf(v.Type) {
			return v, CastException(c.Type, v.Type)
		}
		casted, err := v.Type.Cast(c)
		if err != nil {
			return v, err
		}
		v.Data = casted.Data
		return v, nil
	case Reference:
		return v, ReferenceValueException()
	}
	return v, nil
}

type Reference struct {
	Type     *Datatype
	Name     string
	Link     *Value
	Constant bool
}

func (r Reference) Alias() string {
	return r.Name
}

func (r Reference) SearchSpace() SearchSpace {
	return IdentifierSearchSpace
}

func (r Reference) Update(item SearchItem) (SearchItem, error) {
	if r.SearchSpace() != item.SearchSpace() {
		return r, SearchSpaceException()
	}
	switch c := item.(type) {
	case Value:
		return r, ValueReferenceException()
	case Reference:
		if r.Constant {
			return r, ConstantException(r.Name)
		}
		if !c.Type.KindOf(r.Type) {
			return r, CastException(c.Type, r.Type)
		}
		r.Link = c.Link
	}
	return r, nil
}

type SearchItem interface {
	SearchSpace() SearchSpace
	Alias() string
	Update(SearchItem) (SearchItem, error)
}

type SearchSpace int

const (
	IdentifierSearchSpace SearchSpace = iota
	OperatorSearchSpace
	DatatypeSearchSpace
)

var (
	SearchSpaces = []SearchSpace{
		IdentifierSearchSpace,
		OperatorSearchSpace,
		DatatypeSearchSpace,
	}
)

type Namespace struct {
	Parent  *Namespace
	Storage map[SearchSpace]map[string]SearchItem
}

func (ns *Namespace) Find(space SearchSpace, alias string) (SearchItem, error) {
	searchSpace := ns.Storage[space]
	item, ok := searchSpace[alias]
	if !ok {
		if ns.Parent != nil {
			return ns.Parent.Find(space, alias)
		}
		return nil, NamespaceException(alias)
	}
	return item, nil
}

func (ns *Namespace) Update(item SearchItem) error {
	existing, ok := ns.Storage[item.SearchSpace()][item.Alias()]
	if !ok {
		if ns.Parent == nil {
			return NamespaceException(item.Alias())
		}
		return ns.Parent.Update(item)
	}
	existing, err := existing.Update(item)
	if err != nil {
		return err
	}
	ns.Storage[item.SearchSpace()][item.Alias()] = existing
	return nil
}

func (ns *Namespace) Store(item SearchItem) error {
	ns.Storage[item.SearchSpace()][item.Alias()] = item
	return nil
}

func (ns *Namespace) Child() *Namespace {
	return NewNamespace(ns)
}

func NewNamespace(parent *Namespace) *Namespace {
	storage := make(map[SearchSpace]map[string]SearchItem)
	for _, ss := range SearchSpaces {
		storage[ss] = make(map[string]SearchItem)
	}
	ns := &Namespace{
		Parent:  parent,
		Storage: storage,
	}
	return ns
}

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
