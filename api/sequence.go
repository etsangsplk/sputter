package api

import (
	"bytes"
	"strconv"
)

// IndexNotFound is thrown if an index is not found in a sequence
const IndexNotFound = "index not found: %s"

type (
	// Invoker is the standard signature for a function that can be invoked
	Invoker func(Context, Values) Value

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
func IndexedApply(s Indexed, args Values) Value {
	i := AssertArityRange(args, 1, 2)
	idx := AssertInteger(args[0])
	if r, ok := s.ElementAt(idx); ok {
		return r
	}
	if i == 2 {
		return args[1]
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
	return v.(CountedSequence).Count()
}
