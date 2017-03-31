package api

import "fmt"

const (
	// ExpectedCountable is thrown if taking count of a non-countable sequence
	ExpectedCountable = "sequence is not countable"

	// ExpectedSequence is thrown when a Value is not a Sequence
	ExpectedSequence = "value '%s' is not a sequence"
)

// SequenceProcessor is the standard signature for a function that is
// capable of transforming or validating a Sequence
type SequenceProcessor func(Context, Sequence) Value

// Sequence interfaces expose a lazily resolved sequence of Values
type Sequence interface {
	First() Value
	Rest() Sequence
	Prepend(v Value) Sequence
	IsSequence() bool
}

// Countable interfaces allow a Sequence to return a count of its items
type Countable interface {
	Sequence
	Count() int
}

// Indexed interfaces allow a Sequence item to be retrieved by index
type Indexed interface {
	Sequence
	Get(index int) Value
}

// Iterator is a stateful iteration interface for Sequences.  "Stateful"
// is the key word here.  This data structure should not be used in any
// concurrent or immutable situation
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

// Count will return the Count from a Countable Sequence or explode
func Count(s Sequence) int {
	if f, ok := s.(Countable); ok {
		return f.Count()
	}
	panic(ExpectedCountable)
}

// AssertSequence will cast a Value into a Sequence or explode violently
func AssertSequence(v Value) Sequence {
	if r, ok := v.(Sequence); ok {
		return r
	}
	panic(fmt.Sprintf(ExpectedSequence, String(v)))
}
