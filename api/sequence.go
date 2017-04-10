package api

import "strconv"

const (
	// ExpectedCountable is thrown if taking count of a non-countable sequence
	ExpectedCountable = "sequence is not countable"

	// ExpectedSequence is thrown when a Value is not a Sequence
	ExpectedSequence = "value is not a sequence: %s"

	// ExpectedIndexed is thrown when a Value is not Indexed
	ExpectedIndexed = "value is not an indexed sequence: %s"

	// IndexNotFound is thrown if an index is not found in a sequence
	IndexNotFound = "index not found in sequence: %s"
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
	Get(index int) (Value, bool)
}

// Iterator is a stateful iteration interface for Sequences. "Stateful"
// is the key word here. This data structure should not be used in any
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

// IndexedApply provides 'nth' behavior for Indexed Sequences
func IndexedApply(s Indexed, c Context, args Sequence) Value {
	i := AssertArityRange(args, 1, 2)
	idx := AssertInteger(Eval(c, args.First()))
	if r, ok := s.Get(idx); ok {
		return r
	}
	if i == 2 {
		return Eval(c, args.Rest().First())
	}
	panic(Err(IndexNotFound, strconv.Itoa(idx)))
}

// AssertSequence will cast a Value into a Sequence or explode violently
func AssertSequence(v Value) Sequence {
	if r, ok := v.(Sequence); ok {
		return r
	}
	panic(Err(ExpectedSequence, String(v)))
}

// AssertIndexed will cast a Value into an Indexed or explode violently
func AssertIndexed(v Value) Indexed {
	if r, ok := v.(Indexed); ok {
		return r
	}
	panic(Err(ExpectedIndexed, String(v)))
}
