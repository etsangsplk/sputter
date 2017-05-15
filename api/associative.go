package api

import "bytes"

const (
	// ExpectedPair is thrown if you prepend to a Map incorrectly
	ExpectedPair = "expected two element vectors when prepending"

	// ExpectedMapped is thrown if the Value is not a Mapped item
	ExpectedMapped = "expected a mapped sequence: %s"

	// KeyNotFound is thrown if a Key is not found in a Mapped item
	KeyNotFound = "key not found in mapped sequence: %s"
)

// Mapped interfaces allow a Sequence item to be retrieved by Name
type Mapped interface {
	Sequence
	Getter
}

// Associative is a Mappable that is implemented atop an array
type Associative interface {
	Conjoiner
	MakeEvaluable
	Getter
	Counted
	Applicable
	Associative() bool
}

type associative []Vector

type evaluableAssociative struct {
	associative
}

// NewAssociative instantiates a new Associative
func NewAssociative(v ...Vector) Associative {
	return associative(v)
}

// Associative is a disambiguating marker
func (a associative) Associative() bool {
	return true
}

// Count returns the number of key/value pairs in this Associative
func (a associative) Count() int {
	return len(a)
}

// Get returns the Value corresponding to the key in the Associative
func (a associative) Get(key Value) (Value, bool) {
	l := len(a)
	for i := 0; i < l; i++ {
		mp := a[i]
		k, _ := mp.ElementAt(0)
		if k == key {
			v, _ := mp.ElementAt(1)
			return v, true
		}
	}
	return Nil, false
}

// Apply makes Associative applicable
func (a associative) Apply(c Context, args Sequence) Value {
	AssertArity(args, 1)
	k := Eval(c, args.First())
	if r, ok := a.Get(k); ok {
		return r
	}
	panic(Err(KeyNotFound, k))
}

// First returns the first key/value pair of an Associative
func (a associative) First() Value {
	return a[0]
}

// Rest returns the remaining elements of an Associative as a Sequence
func (a associative) Rest() Sequence {
	return Sequence(a[1:])
}

// Prepend creates a new Sequence by prepending a Value
func (a associative) Prepend(v Value) Sequence {
	if mp, ok := v.(Vector); ok {
		AssertArity(mp, 2)
		return append(associative{mp}, a...)
	}
	panic(ExpectedPair)
}

// Conjoin implements the Conjoiner interface
func (a associative) Conjoin(v Value) Sequence {
	return a.Prepend(v)
}

// IsSequence returns whether this instance is a consumable Sequence
func (a associative) IsSequence() bool {
	return len(a) > 0
}

// Evaluable turns Associative into an Evaluable Expression
func (a associative) Evaluable() Value {
	return &evaluableAssociative{
		associative: a,
	}
}

// Str converts this Value into a Str
func (a associative) Str() Str {
	var b bytes.Buffer
	l := len(a)

	b.WriteString("{")
	for i := 0; i < l; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		mp := a[i]
		k, _ := mp.ElementAt(0)
		v, _ := mp.ElementAt(1)
		b.WriteString(string(k.Str()))
		b.WriteString(" ")
		b.WriteString(string(v.Str()))
	}
	b.WriteString("}")
	return Str(b.String())
}

// Eval makes an EvaluableAssociative Evaluable
func (e *evaluableAssociative) Eval(c Context) Value {
	a := e.associative
	l := len(a)
	r := make(associative, l)
	for i := 0; i < l; i++ {
		mp := a[i]
		k, _ := mp.ElementAt(0)
		v, _ := mp.ElementAt(1)
		r[i] = NewVector(Eval(c, k), Eval(c, v))
	}
	return r
}

// AssertMapped will cast Value to a Mapped or explode violently
func AssertMapped(v Value) Mapped {
	if r, ok := v.(Mapped); ok {
		return r
	}
	panic(Err(ExpectedMapped, v))
}
