package api

import (
	"bytes"
)

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

// Iterate creates a new Iterator instance for the Vector
func (v Vector) Iterate() Iterator {
	return &vectorIterator{v, len(v), 0}
}

// Evaluate makes a Vector Evaluable
func (v Vector) Evaluate(c *Context) Value {
	r := make(Vector, len(v))
	for i := 0; i < len(v); i++ {
		r[i] = Evaluate(c, v[i])
	}
	return r
}

type vectorIterator struct {
	vector Vector
	len    int
	pos    int
}

// Next returns the next Value from the Iterator
func (i *vectorIterator) Next() (v Value, ok bool) {
	if i.pos < i.len {
		r := i.vector[i.pos]
		i.pos++
		return r, true
	}
	return Nil, false
}

// Iterable returns a new Iterable from the Iterator's current state
func (i *vectorIterator) Iterable() Iterable {
	return i.vector[i.pos:]
}

func (v Vector) String() string {
	var b bytes.Buffer

	b.WriteString("[")
	for i := 0; i < len(v); i++ {
		vi := v[i]
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(ValueToString(vi))
	}
	b.WriteString("]")
	return b.String()
}
