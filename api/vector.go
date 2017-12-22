package api

import "bytes"

type (
	// Vector is a fixed-length array of Values
	Vector interface {
		Conjoiner
		Indexed
		Counted
		Applicable
		Evaluable
		VectorType()
	}

	// Values is the concrete implementation of a Vector
	Values []Value
)

var emptyVector = Values{}

// NewVector instantiates a new Vector
func NewVector(v ...Value) Vector {
	return Values(v)
}

// VectorType identifies this Value as a Vector
func (Values) VectorType() {}

// Count returns the number of elements in the Value array
func (v Values) Count() int {
	return len(v)
}

// ElementAt returns a specific element of the Value array
func (v Values) ElementAt(index int) (Value, bool) {
	vals := v
	if index >= 0 && index < len(vals) {
		return vals[index], true
	}
	return Nil, false
}

// Apply makes Values Applicable
func (v Values) Apply(_ Context, args Sequence) Value {
	return IndexedApply(v, args)
}

// Eval evaluates its elements, returning a new Value array
func (v Values) Eval(c Context) Value {
	l := len(v)
	r := make(Values, l)
	for i := 0; i < l; i++ {
		r[i] = Eval(c, v[i])
	}
	return r
}

// First returns the first element of the Value array
func (v Values) First() Value {
	if len(v) > 0 {
		return v[0]
	}
	return Nil
}

// Rest returns the elements of the Value array that follow the first
func (v Values) Rest() Sequence {
	if len(v) > 1 {
		return v[1:]
	}
	return emptyVector
}

// IsSequence returns whether or not this Value array has any elements
func (v Values) IsSequence() bool {
	return len(v) > 0
}

// Split breaks the Value array into its components (first, rest, isSequence)
func (v Values) Split() (Value, Sequence, bool) {
	lv := len(v)
	if lv > 1 {
		return v[0], v[1:], true
	} else if lv == 1 {
		return v[0], emptyVector, true
	}
	return Nil, emptyVector, false
}

// Prepend inserts an element at the beginning of the Value array
func (v Values) Prepend(p Value) Sequence {
	return append(Values{p}, v...)
}

// Conjoin appends an element to the end of the Value array
func (v Values) Conjoin(a Value) Sequence {
	return append(v, a)
}

// Concat concatenates two Value arrays
func (v Values) Concat(a Values) Values {
	return append(v, SequenceToValues(a)...)
}

// Str converts this Value array to a Str
func (v Values) Str() Str {
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
