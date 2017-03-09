package api

// Sequence interfaces expose a lazily resolved sequence of Values
type Sequence interface {
	First() Value
	Rest() Sequence
	Prepend(v Value) Sequence
	IsSequence() bool
}

// Countable interfaces allow a Sequence to return a count of its items
type Countable interface {
	Count() int
}

// Indexed interfaces allow a Sequence item to be retrieved by index
type Indexed interface {
	Get(index int) Value
}

// Mapped interfaces allow a Sequence item to be retrieved by Name
type Mapped interface {
	Get(key Value) Value
}

// Iterator is a stateful iteration interface for Sequences
type Iterator struct {
	sequence Sequence
}

// Next returns the next value from the Iterator
func (i *Iterator) Next() (Value, bool) {
	s := i.sequence
	if !s.IsSequence() {
		return Nil, false
	}
	r := s.First()
	i.sequence = s.Rest()
	return r, true
}

// Rest returns the rest of the Iteration as a Sequence
func (i *Iterator) Rest() Sequence {
	return i.sequence
}

// Iterate creates a stateful Iterator over a Sequence
func Iterate(s Sequence) *Iterator {
	return &Iterator{sequence: s}
}
