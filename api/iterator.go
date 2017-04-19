package api

// Iterator is a stateful iteration interface for Sequences. "Stateful"
// is the key word here. This data structure should not be used in any
// concurrent or immutable situation.
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
