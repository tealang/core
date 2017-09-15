package runtime

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
	Type     *Datatype
	Data     interface{}
	Name     string
	Constant bool
}

func (v Value) String() string {
	if v.Type == nil {
		return "null"
	}
	return v.Type.Format(v)
}

// Alias returns the name of the value.
func (v Value) Alias() string {
	return v.Name
}

// SearchSpace returns the search space of the value.
func (v Value) SearchSpace() SearchSpace {
	return IdentifierSearchSpace
}

// Update sets the data of the value.
func (v Value) Update(item SearchItem) (SearchItem, error) {
	if v.SearchSpace() != item.SearchSpace() {
		return v, SearchSpaceException{}
	}
	switch c := item.(type) {
	case Value:
		if v.Constant {
			return v, ConstantException{v.Name}
		}
		if !c.Type.KindOf(v.Type) {
			return v, CastException{
				From: c.Type,
				To:   v.Type,
			}
		}
		casted, err := v.Type.Cast(c)
		if err != nil {
			return v, err
		}
		v.Data = casted.Data
		return v, nil
	case Reference:
		return v, ReferenceValueException{}
	}
	return v, nil
}

// Reference of a given datatype stores a link to the value.
type Reference struct {
	Type     *Datatype
	Name     string
	Link     *Value
	Constant bool
}

// Alias returns the reference name.
func (r Reference) Alias() string {
	return r.Name
}

// SearchSpace returns the reference search space.
func (r Reference) SearchSpace() SearchSpace {
	return IdentifierSearchSpace
}

// Update sets the reference link.
func (r Reference) Update(item SearchItem) (SearchItem, error) {
	if r.SearchSpace() != item.SearchSpace() {
		return r, SearchSpaceException{}
	}
	switch c := item.(type) {
	case Value:
		return r, ValueReferenceException{}
	case Reference:
		if r.Constant {
			return r, ConstantException{r.Name}
		}
		if !c.Type.KindOf(r.Type) {
			return r, CastException{
				From: c.Type,
				To:   r.Type,
			}
		}
		r.Link = c.Link
	}
	return r, nil
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
	// IdentifierSearchSpace is used for variables and functions.
	IdentifierSearchSpace SearchSpace = iota
	// OperatorSearchSpace is used for operators.
	OperatorSearchSpace
	// DatatypeSearchSpace is used for datatypes.
	DatatypeSearchSpace
)

var (
	// SearchSpaces lists all available search spaces.
	SearchSpaces = []SearchSpace{
		IdentifierSearchSpace,
		OperatorSearchSpace,
		DatatypeSearchSpace,
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
