package runtime

import "fmt"

// Formatter formats the given value in the datatype format.
type Formatter func(v Value) string

// Caster attempts to convert the value in a value conforming to a specific datatype.
type Caster func(v Value) (Value, error)

// Datatype has a name, a parent, a caster and a formatter.
type Datatype struct {
	Name   string
	Parent *Datatype
	Cast   Caster
	Format Formatter
}

func (Datatype) SearchSpace() SearchSpace {
	return SearchDatatype
}

func (datatype *Datatype) Alias() string {
	return datatype.Name
}

func (datatype *Datatype) Update(item SearchItem) (SearchItem, error) {
	return item, Exception{"Datatypes cannot be overriden"}
}

// KindOf checks if this datatype is of the same kind as the given datatype.
func (datatype *Datatype) KindOf(other *Datatype) bool {
	if datatype != other {
		if datatype.Parent != nil {
			return datatype.Parent.KindOf(other)
		}
		return false
	}
	return true
}

func (datatype *Datatype) String() string {
	return datatype.Name
}

// Value of a given datatype stored data associated with a name.
type Value struct {
	Type      *Datatype
	Data      interface{}
	Name      string
	Constant  bool
	Reference bool
}

func (v Value) Rechange(constant bool) Value {
	return Value{
		Type:      v.Type,
		Data:      v.Data,
		Name:      v.Name,
		Constant:  constant,
		Reference: v.Reference,
	}
}

func (v Value) Rename(alias string) Value {
	return Value{
		Type:      v.Type,
		Data:      v.Data,
		Name:      alias,
		Constant:  v.Constant,
		Reference: v.Reference,
	}
}

func (v Value) VariableString() string {
	if v.Type == nil {
		return ":null"
	}

	if v.Reference {
		return fmt.Sprintf("%s: &%s", v.Name, v.Type)
	}
	return fmt.Sprintf("%s: %s", v.Name, v.Type)
}

func (v Value) String() string {
	if v.Type == nil {
		return "null"
	}

	if v.Reference {
		return v.Type.Format(*v.Data.(*Value))
	}
	return v.Type.Format(v)
}

// Alias returns the name of the value.
func (v Value) Alias() string {
	return v.Name
}

// SearchSpace returns the search space of the value.
func (v Value) SearchSpace() SearchSpace {
	return SearchIdentifier
}

// Update sets the data of the value.
func (v Value) Update(item SearchItem) (SearchItem, error) {
	if v.Constant {
		return v, ConstantException{v.Name}
	}
	if v.SearchSpace() != item.SearchSpace() {
		return v, SearchSpaceException{}
	}
	c, ok := item.(Value)
	if !ok {
		return v, UnexpectedItemException{Expected: Value{}, Got: item}
	}
	if !c.Type.KindOf(v.Type) {
		return v, CastException{
			From: c.Type,
			To:   v.Type,
		}
	}

	if v.Reference {
		if !c.Reference {
			return v, ReferenceValueException{}
		}
		v.Data = c.Data
	} else {
		casted, err := v.Type.Cast(c)
		if err != nil {
			return v, err
		}
		v.Data = casted.Data
	}

	return v, nil
}

type Operator struct {
	Function
	Symbol   string
	Constant bool
}

func (Operator) SearchSpace() SearchSpace {
	return SearchOperator
}

func (o Operator) Alias() string {
	return o.Symbol
}

func (o Operator) Update(item SearchItem) (SearchItem, error) {
	if o.Constant {
		return o, ConstantException{o.Symbol}
	}
	if o.SearchSpace() != item.SearchSpace() {
		return o, SearchSpaceException{}
	}
	op, ok := item.(Operator)
	if !ok {
		return o, UnexpectedItemException{Operator{}, item}
	}
	o.Signatures = op.Signatures
	return o, nil
}

// SearchItem is a generic item that can be stored in a namespace.
type SearchItem interface {
	SearchSpace() SearchSpace
	Alias() string
	Update(SearchItem) (SearchItem, error)
}

// SearchSpace identifiers the space, where a search item can be stored and found.
type SearchSpace int

const (
	// SearchIdentifier is used for variables and functions.
	SearchIdentifier SearchSpace = iota
	// SearchOperator is used for operators.
	SearchOperator
	// SearchDatatype is used for datatypes.
	SearchDatatype
)

var (
	// SearchSpaces lists all available search spaces.
	SearchSpaces = []SearchSpace{
		SearchIdentifier,
		SearchOperator,
		SearchDatatype,
	}
)

// Namespace stores search items in a hierarchical structure.
type Namespace struct {
	Parent  *Namespace
	Storage map[SearchSpace]map[string]SearchItem
}

// Find looks for the search item in the given search space of this namespace and its parent.
func (ns *Namespace) Find(space SearchSpace, alias string) (SearchItem, error) {
	searchSpace := ns.Storage[space]
	item, ok := searchSpace[alias]
	if !ok {
		if ns.Parent != nil {
			return ns.Parent.Find(space, alias)
		}
		return nil, NamespaceException{alias}
	}
	return item, nil
}

// Update looks for and updates the search item in this namespace and its parent.
func (ns *Namespace) Update(item SearchItem) error {
	existing, ok := ns.Storage[item.SearchSpace()][item.Alias()]
	if !ok {
		if ns.Parent == nil {
			return NamespaceException{item.Alias()}
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

// Store puts the search item in a search space in this namespace.
func (ns *Namespace) Store(item SearchItem) error {
	_, ok := ns.Storage[item.SearchSpace()][item.Alias()]
	if ok {
		return StoreExistsException{item.Alias()}
	}
	ns.Storage[item.SearchSpace()][item.Alias()] = item
	return nil
}

// Child returns a new namespace that has this namespace as its parent.
func (ns *Namespace) Child() *Namespace {
	return NewNamespace(ns)
}

// NewNamespace initializes a new empty namespace.
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
