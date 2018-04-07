package api

import "bytes"

// Vector is a fixed-length array of Vector
type Vector []Value

// EmptyVector is an empty Vector
var EmptyVector = Vector{}

// NewVector instantiates a new Vector
func NewVector(v ...Value) Vector {
	return Vector(v)
}

// Count returns the number of elements in the Vector
func (v Vector) Count() int {
	return len(v)
}

// ElementAt returns a specific element of the Vector
func (v Vector) ElementAt(index int) (Value, bool) {
	vals := v
	if index >= 0 && index < len(vals) {
		return vals[index], true
	}
	return Nil, false
}

// Apply makes Vector Applicable
func (v Vector) Apply(_ Context, args Vector) Value {
	return IndexedApply(v, args)
}

// Eval evaluates its elements, returning a new Vector
func (v Vector) Eval(c Context) Value {
	l := len(v)
	r := make(Vector, l)
	for i := 0; i < l; i++ {
		r[i] = Eval(c, v[i])
	}
	return r
}

// First returns the first element of the Vector
func (v Vector) First() Value {
	if len(v) > 0 {
		return v[0]
	}
	return Nil
}

// Rest returns the elements of the Vector that follow the first
func (v Vector) Rest() Sequence {
	if len(v) > 1 {
		return v[1:]
	}
	return EmptyVector
}

// IsSequence returns whether or not this Vector has any elements
func (v Vector) IsSequence() bool {
	return len(v) > 0
}

// Split breaks the Vector into its components (first, rest, isSequence)
func (v Vector) Split() (Value, Sequence, bool) {
	lv := len(v)
	if lv > 1 {
		return v[0], v[1:], true
	} else if lv == 1 {
		return v[0], EmptyVector, true
	}
	return Nil, EmptyVector, false
}

// Prepend inserts an element at the beginning of the Vector
func (v Vector) Prepend(p Value) Sequence {
	return append(Vector{p}, v...)
}

// Conjoin appends an element to the end of the Vector
func (v Vector) Conjoin(a Value) Sequence {
	return append(v, a)
}

// Concat concatenates two Value arrays
func (v Vector) Concat(a Vector) Vector {
	return append(v, SequenceToVector(a)...)
}

// Str converts this Vector to a Str
func (v Vector) Str() Str {
	var b bytes.Buffer
	l := len(v)

	b.WriteString("[")
	for i := 0; i < l; i++ {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(string(v[i].Str()))
	}
	b.WriteString("]")
	return Str(b.String())
}
