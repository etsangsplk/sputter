package api

import (
	"bytes"
	"fmt"
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
	result := make(Vector, len(v))
	for index := 0; index < len(v); index++ {
		result[index] = Evaluate(c, v[index])
	}
	return result
}

type vectorIterator struct {
	vector Vector
	len    int
	pos    int
}

// Next returns the next Value from the Iterator
func (i *vectorIterator) Next() (Value, bool) {
	if i.pos < i.len {
		result := i.vector[i.pos]
		i.pos++
		return result, true
	}
	return EmptyList, false
}

// Iterable returns a new Iterable from the Iterator's current state
func (i *vectorIterator) Iterable() Iterable {
	return i.vector[i.pos:]
}

func (v Vector) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	for i := 0; i < len(v); i++ {
		current := v[i]
		if i > 0 {
			buffer.WriteString(" ")
		}
		if str, ok := current.(fmt.Stringer); ok {
			buffer.WriteString(str.String())
		} else {
			buffer.WriteString(current.(string))
		}
	}
	buffer.WriteString("]")
	return buffer.String()
}
