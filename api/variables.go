package api

// Comparison represents the result of a equality comparison
type Comparison int

const (
	// LessThan means left Value is less than right Value
	LessThan Comparison = -1

	// EqualTo means left Value is equal to right Value
	EqualTo Comparison = 0

	// GreaterThan means left Value is greater than right Value
	GreaterThan Comparison = 1
)

// Value is the generic interface for all 'Values'
type Value interface {
	Str() Str
}

// Name is a Variable name
type Name string

// Bool represents the values True or False
type Bool bool

type atom struct {
	str Str
}

// Comparer is an interface for a Value capable of comparing.
type Comparer interface {
	Compare(Comparer) Comparison
}

// Variables represents a mapping from Name to Value
type Variables map[Name]Value

// Typed is the generic interface for Values that are typed
type Typed interface {
	Type() Name
}

// Getter is the interface for Values that have retrievable Properties
type Getter interface {
	Get(key Value) (Value, bool)
}

// Elementer is the interface for Values that have indexed elements
type Elementer interface {
	ElementAt(index int) (Value, bool)
}

// Named is the generic interface for Values that are named
type Named interface {
	Name() Name
}

// Name makes Name Named
func (n Name) Name() Name {
	return n
}

// Str converts this Value into a Str
func (n Name) Str() Str {
	return Str(n)
}

// Str converts this Value into a Str
func (b Bool) Str() Str {
	if bool(b) {
		return "true"
	}
	return "false"
}

// Atom instantiates a new Atom
func Atom(str Str) Value {
	return &atom{str: str}
}

// Str converts this Value into a Str
func (a *atom) Str() Str {
	return a.str
}
