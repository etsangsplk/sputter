package sputter

// Value is the generic interface for all 'Values' in the VM
type Value interface{
}

// Iterable values can be used in loops and comprehensions
type Iterable interface {
	Iterate() Iterator
}

// Iterator functions are ones that return the next Value of an Iterable
type Iterator func() (Value, bool)

// Literal identifies a Value as being a literal reference
type Literal struct {
	value Value
}

// Resolvable is an Identifier that can be resolved against a Context
type Resolvable struct {
	name string
}
