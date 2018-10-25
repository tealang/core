package runtime

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Formatter formats the given value in the datatype format.
type Formatter func(v Value) string

// Caster attempts to convert the value in a value conforming to a specific datatype.
type Caster func(v Value, f []Typeflag) (Value, error)

// Datatype has a name, a parent, a caster and a formatter.
type Datatype struct {
	Name   string
	Parent *Datatype
	Cast   Caster
	Format Formatter
}

// SearchSpace returns the Datatype search space.
func (Datatype) SearchSpace() SearchSpace {
	return SearchDatatype
}

// Alias returns the datatypes name.
func (datatype *Datatype) Alias() string {
	return datatype.Name
}

// Update fails.
func (datatype *Datatype) Update(item SearchItem) (SearchItem, error) {
	return item, errors.Errorf("datatype %s can not be updated", datatype.Name)
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

// Typeflag of a given datatype with one or more children typeflags.
type Typeflag struct {
	Type   *Datatype
	Params []Typeflag
}

func (tf Typeflag) String() string {
	if len(tf.Params) > 0 {
		params := make([]string, len(tf.Params))
		for i := range tf.Params {
			params[i] = tf.Params[i].String()
		}
		return fmt.Sprintf("%s<%s>", tf.Type.Name, strings.Join(params, ", "))
	}
	return tf.Type.Name
}

// Cast does a cast to the type of the source value including source type params.
func (tf Typeflag) Cast(v Value) (Value, error) {
	return tf.Type.Cast(v, tf.Params)
}

// T builds a typeflag using the given list of types.
// For example, T(array, any, string) -> array<any<string>>
func T(tree ...*Datatype) Typeflag {
	if len(tree) == 0 {
		return Typeflag{}
	} else if len(tree) == 1 {
		return Typeflag{
			Type: tree[0],
		}
	} else {
		return Typeflag{
			Type:   tree[0],
			Params: []Typeflag{T(tree[1:]...)},
		}
	}
}

// Value of a given datatype stored data associated with a name.
type Value struct {
	Typeflag
	Data      interface{}
	Name      string
	Constant  bool
	Reference bool
}

// EqualTo checks if the both values are equal.
// To be classified as equal, they must have the same type and data.
func (v Value) EqualTo(w Value) bool {
	return v.Type == w.Type && v.Data == w.Data
}

// Rechange turns a (constant/variable) value into a (constant/variable) value.
func (v Value) Rechange(constant bool) Value {
	return Value{
		Typeflag:  v.Typeflag,
		Data:      v.Data,
		Name:      v.Name,
		Constant:  constant,
		Reference: v.Reference,
	}
}

// Rename changes the values name.
func (v Value) Rename(alias string) Value {
	return Value{
		Typeflag:  v.Typeflag,
		Data:      v.Data,
		Name:      alias,
		Constant:  v.Constant,
		Reference: v.Reference,
	}
}

// VariableString generates a representation in the variable declaration format.
func (v Value) VariableString() string {
	if v.Type == nil {
		return ":null"
	}
	if v.Reference {
		return fmt.Sprintf("%s: &%s", v.Name, v.Typeflag)
	}
	return fmt.Sprintf("%s: %s", v.Name, v.Typeflag)
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
		return v, errors.Errorf("value %s can not be changed", v)
	}
	if v.SearchSpace() != item.SearchSpace() {
		return v, errors.Errorf("item to update is from different search space")
	}
	c, ok := item.(Value)
	if !ok {
		return v, errors.Errorf("expected value item, got %s", item)
	}
	if !c.Type.KindOf(v.Type) {
		return v, errors.Errorf("can not assign type %s to %s", c.Type, v.Type)
	}

	if v.Reference {
		if !c.Reference {
			return v, errors.Errorf("value can not be assigned to reference %s", v.Name)
		}
		v.Data = c.Data
	} else {
		casted, err := v.Cast(c)
		if err != nil {
			return v, errors.Wrap(err, "could not update value")
		}
		v.Data = casted.Data
	}

	return v, nil
}

// Operator is a function with an associated operation symbol.
type Operator struct {
	Function
	Symbol   string
	Constant bool
}

// SearchSpace returns the operator search space.
func (Operator) SearchSpace() SearchSpace {
	return SearchOperator
}

// Alias returns the operators symbol.
func (o Operator) Alias() string {
	return o.Symbol
}

// Update changes the function and state of the operator.
func (o Operator) Update(item SearchItem) (SearchItem, error) {
	if o.Constant {
		return o, errors.Errorf("operator %s can not be changed", o.Symbol)
	}
	if o.SearchSpace() != item.SearchSpace() {
		return o, errors.Errorf("item to update is from different search space")
	}
	op, ok := item.(Operator)
	if !ok {
		return o, errors.Errorf("expected item operator, got %s", item)
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
		return nil, errors.Errorf("item %s not found in namespace", alias)
	}
	return item, nil
}

// Update looks for and updates the search item in this namespace and its parent.
func (ns *Namespace) Update(item SearchItem) error {
	existing, ok := ns.Storage[item.SearchSpace()][item.Alias()]
	if !ok {
		if ns.Parent == nil {
			return errors.Errorf("item %s not found in namespace", item.Alias())
		}
		return ns.Parent.Update(item)
	}
	existing, err := existing.Update(item)
	if err != nil {
		return errors.Wrap(err, "update failed")
	}
	ns.Storage[item.SearchSpace()][item.Alias()] = existing
	return nil
}

// Store puts the search item in a search space in this namespace.
func (ns *Namespace) Store(item SearchItem) error {
	_, ok := ns.Storage[item.SearchSpace()][item.Alias()]
	if ok {
		return errors.Errorf("item %s already exists in namespace", item.Alias())
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
