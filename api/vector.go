package api

import "bytes"

// ExpectedVector is raised if a value is not a Vector
const ExpectedVector = "value is not a vector: %s"

// Vector is a fixed-length Array of Values
type Vector interface {
	Conjoiner
	MakeEvaluable
	Elementer
	Counted
	Applicable
	Vector() bool
}

type vector []Value

type evaluableVector struct {
	vector
}

var emptyVector = vector{}

// NewVector instantiates a new Vector
func NewVector(v ...Value) Vector {
	return vector(v)
}

// Vector is a disambiguating marker
func (v vector) Vector() bool {
	return true
}

// Count returns the length of the Vector
func (v vector) Count() int {
	return len(v)
}

// ElementAt returns the Value at the indexed position in the Vector
func (v vector) ElementAt(index int) (Value, bool) {
	if index >= 0 && index < len(v) {
		return v[index], true
	}
	return Nil, false
}

// Apply makes Vector applicable
func (v vector) Apply(c Context, args Sequence) Value {
	return IndexedApply(v, c, args)
}

// First returns the first element of a Vector
func (v vector) First() Value {
	if len(v) > 0 {
		return v[0]
	}
	return Nil
}

// Rest returns the remaining elements of a Vector as a Sequence
func (v vector) Rest() Sequence {
	if len(v) > 1 {
		return Sequence(v[1:])
	}
	return emptyVector
}

// Prepend creates a new Sequence by prepending a Value
func (v vector) Prepend(p Value) Sequence {
	return append(vector{p}, v...)
}

// Conjoin implements the Conjoiner interface
func (v vector) Conjoin(a Value) Sequence {
	return append(v, a)
}

// IsSequence returns whether this instance is a consumable Sequence
func (v vector) IsSequence() bool {
	return len(v) > 0
}

// Evaluable turns Vector into an Evaluable Expression
func (v vector) Evaluable() Value {
	return &evaluableVector{
		vector: v,
	}
}

// Str converts this Value into a Str
func (v vector) Str() Str {
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

// Eval makes an evaluableVector Evaluable
func (e *evaluableVector) Eval(c Context) Value {
	v := e.vector
	l := len(v)
	r := make(vector, l)
	for i := 0; i < l; i++ {
		r[i] = Eval(c, v[i])
	}
	return r
}

// AssertVector will cast the Value into a Vector or die trying
func AssertVector(v Value) Vector {
	if r, ok := v.(Vector); ok {
		return r
	}
	panic(Err(ExpectedVector, v))
}
