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
	ExpectedIndexed = "expected an indexed value: %s"

	// ExpectedConjoiner is thrown when a Value is not a Conjoiner
	ExpectedConjoiner = "value can not be conjoined: %s"

	// IndexNotFound is thrown if an index is not found in a sequence
	IndexNotFound = "index not found: %s"
)

type (
	// SequenceProcessor is the standard signature for a ReflectedFunction that is
	// capable of transforming or validating a Sequence
	SequenceProcessor func(Context, Sequence) Value

	// Sequence interfaces expose a lazily resolved sequence of Values
	Sequence interface {
		Value
		First() Value
		Rest() Sequence
		Split() (Value, Sequence, bool)
		Prepend(Value) Sequence
		IsSequence() bool
	}

	// Conjoiner is a Sequence that can be Conjoined in some way
	Conjoiner interface {
		Sequence
		Conjoin(Value) Sequence
	}

	// IndexedSequence is a Sequence that provides an Indexed interface
	IndexedSequence interface {
		Sequence
		Indexed
	}

	// CountedSequence is a Sequence that provides a Counted interface
	CountedSequence interface {
		Sequence
		Counted
	}

	// MappedSequence is a Sequence that provides a Mapped interface
	MappedSequence interface {
		Sequence
		Mapped
	}
)

// IndexedApply provides 'nth' behavior for Indexed Sequences
func IndexedApply(s Indexed, args Sequence) Value {
	i := AssertArityRange(args, 1, 2)
	idx := AssertInteger(args.First())
	if r, ok := s.ElementAt(idx); ok {
		return r
	}
	if i == 2 {
		return args.Rest().First()
	}
	panic(ErrStr(IndexNotFound, strconv.Itoa(idx)))
}

// MakeSequenceStr converts a Sequence to a Str
func MakeSequenceStr(s Sequence) Str {
	f, r, ok := s.Split()
	if !ok {
		return "()"
	}

	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(string(f.Str()))
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		b.WriteString(" ")
		b.WriteString(string(f.Str()))
	}
	b.WriteString(")")
	return Str(b.String())
}

// Count will return the Count from a Counted Sequence or explode
func Count(v Value) int {
	if c, ok := v.(CountedSequence); ok {
		return c.Count()
	}
	panic(ErrStr(ExpectedCounted, v))
}

// AssertSequence will cast a Value into a Sequence or explode violently
func AssertSequence(v Value) Sequence {
	if r, ok := v.(Sequence); ok {
		return r
	}
	panic(ErrStr(ExpectedSequence, v))
}

// AssertIndexed will cast a Value into an Indexed or explode violently
func AssertIndexed(v Value) Indexed {
	if r, ok := v.(Indexed); ok {
		return r
	}
	panic(ErrStr(ExpectedIndexed, v))
}

// AssertConjoiner will cast a Value into a Conjoiner or explode violently
func AssertConjoiner(v Value) Conjoiner {
	if r, ok := v.(Conjoiner); ok {
		return r
	}
	panic(ErrStr(ExpectedConjoiner, v))
}
