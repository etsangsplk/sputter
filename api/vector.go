package api

import "bytes"

// ExpectedVector is raised if a value is not a Vector
const ExpectedVector = "value is not a vector: %s"

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

// Apply makes Vector applicable
func (v Vector) Apply(c Context, args Sequence) Value {
	AssertArity(args, 1)
	idx := AssertInteger(Eval(c, args.First()))
	return v.Get(int(idx))
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

// AssertVector will cast the Value into a Vector or die trying
func AssertVector(v Value) Vector {
	if r, ok := v.(Vector); ok {
		return r
	}
	panic(Err(ExpectedVector, String(v)))
}
