package api

import (
	"bytes"
	"strconv"
)

const (
	// ExpectedCounted is thrown if taking count of a non-countable sequence
	ExpectedCounted = "sequence can not be counted: %s"

	// ExpectedSequence is thrown when a Value is not a Sequence
	ExpectedSequence = "value is not a sequence: %s"

	// ExpectedIndexed is thrown when a Value is not Indexed
	ExpectedIndexed = "value is not an indexed sequence: %s"

	// ExpectedConjoiner is thrown when a Value is not a Conjoiner
	ExpectedConjoiner = "value can not be conjoined: %s"

	// IndexNotFound is thrown if an index is not found in a sequence
	IndexNotFound = "index not found in sequence: %s"
)

// SequenceProcessor is the standard signature for a function that is
// capable of transforming or validating a Sequence
type SequenceProcessor func(Context, Sequence) Value

// Sequence interfaces expose a lazily resolved sequence of Values
type Sequence interface {
	Value
	First() Value
	Rest() Sequence
	Prepend(Value) Sequence
	IsSequence() bool
}

// Counted interfaces allow a Sequence to return a count of its items
type Counted interface {
	Sequence
	Count() int
}

// Indexed interfaces allow a Sequence item to be retrieved by index
type Indexed interface {
	Sequence
	Elementer
}

// Conjoiner interfaces allow a Sequence to be added to
type Conjoiner interface {
	Sequence
	Conjoin(Value) Sequence
}

// Count will return the Count from a Counted Sequence or explode
func Count(v Value) int {
	if c, ok := v.(Counted); ok {
		return c.Count()
	}
	panic(Err(ExpectedCounted, v))
}

// IndexedApply provides 'nth' behavior for Indexed Sequences
func IndexedApply(s Indexed, c Context, args Sequence) Value {
	i := AssertArityRange(args, 1, 2)
	idx := AssertInteger(Eval(c, args.First()))
	if r, ok := s.ElementAt(idx); ok {
		return r
	}
	if i == 2 {
		return Eval(c, args.Rest().First())
	}
	panic(Err(IndexNotFound, strconv.Itoa(idx)))
}

// MakeSequenceStr converts a Sequence to a Str
func MakeSequenceStr(s Sequence) Str {
	if !s.IsSequence() {
		return "()"
	}

	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(string(s.First().Str()))
	for i := s.Rest(); i.IsSequence(); i = i.Rest() {
		b.WriteString(" ")
		b.WriteString(string(i.First().Str()))
	}
	b.WriteString(")")
	return Str(b.String())
}

// AssertSequence will cast a Value into a Sequence or explode violently
func AssertSequence(v Value) Sequence {
	if r, ok := v.(Sequence); ok {
		return r
	}
	panic(Err(ExpectedSequence, v))
}

// AssertIndexed will cast a Value into an Indexed or explode violently
func AssertIndexed(v Value) Indexed {
	if r, ok := v.(Indexed); ok {
		return r
	}
	panic(Err(ExpectedIndexed, v))
}

// AssertConjoiner will cast a Value into a Conjoiner or explode violently
func AssertConjoiner(v Value) Conjoiner {
	if r, ok := v.(Conjoiner); ok {
		return r
	}
	panic(Err(ExpectedConjoiner, v))
}
