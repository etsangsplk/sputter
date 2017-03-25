package api

import "bytes"

// Vector is a fixed-length Array of Values
type Vector []Value

// Count returns the length of the Vector
func (v Vector) Count() int {
	return len(v)
}

// Get returns the Value at the indexed position in the Vector
func (v Vector) Get(index int) Value {
	return v[index]
}

// Eval makes a Vector Evaluable
func (v Vector) Eval(c Context) Value {
	l := len(v)
	r := make(Vector, l)
	for i := 0; i < l; i++ {
		r[i] = Eval(c, v[i])
	}
	return r
}

// First returns the first element of a Vector
func (v Vector) First() Value {
	return v[0]
}

// Rest returns the remaining elements of a Vector as a Sequence
func (v Vector) Rest() Sequence {
	return Sequence(v[1:])
}

// Prepend creates a new Sequence by prepending a Value
func (v Vector) Prepend(p Value) Sequence {
	return append(Vector{p}, v...)
}

// IsSequence returns whether this instance is a consumable Sequence
func (v Vector) IsSequence() bool {
	return len(v) > 0
}

func (v Vector) String() string {
	var b bytes.Buffer
	l := len(v)

	b.WriteString("[")
	for i := 0; i < l; i++ {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(String(v[i]))
	}
	b.WriteString("]")
	return b.String()
}
