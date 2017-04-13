package api

// Name is a Variable name
type Name string

// Value is the generic interface for all 'Values'
type Value interface {
}

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

// Comparer is an interface for a Value capable of comparing.
type Comparer interface {
	Compare(Comparer) Comparison
}

// Named is the generic interface for Values that are named
type Named interface {
	Name() Name
}

// Typed is the generic interface for Values that are typed
type Typed interface {
	Type() Name
}

// Variables represents a mapping from Name to Value
type Variables map[Name]Value

// Name makes Name Named
func (n Name) Name() Name {
	return n
}
