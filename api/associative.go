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
type Associative []Vector

// Count returns the number of key/value pairs in this Associative
func (a Associative) Count() int {
	return len(a)
}

// Get returns the Value corresponding to the key in the Associative
func (a Associative) Get(key Value) (Value, bool) {
	l := len(a)
	for i := 0; i < l; i++ {
		mp := a[i]
		if mp[0] == key {
			return mp[1], true
		}
	}
	return Nil, false
}

// Apply makes Associative applicable
func (a Associative) Apply(c Context, args Sequence) Value {
	AssertArity(args, 1)
	k := Eval(c, args.First())
	if r, ok := a.Get(k); ok {
		return r
	}
	panic(Err(KeyNotFound, k))
}

// Eval makes an Associative Evaluable
func (a Associative) Eval(c Context) Value {
	l := len(a)
	r := make(Associative, l)
	for i := 0; i < l; i++ {
		mp := a[i]
		r[i] = Vector{
			Eval(c, mp[0]),
			Eval(c, mp[1]),
		}
	}
	return r
}

// First returns the first key/value pair of an Associative
func (a Associative) First() Value {
	return a[0]
}

// Rest returns the remaining elements of an Associative as a Sequence
func (a Associative) Rest() Sequence {
	return Sequence(a[1:])
}

// Prepend creates a new Sequence by prepending a Value
func (a Associative) Prepend(v Value) Sequence {
	if mp, ok := v.(Vector); ok {
		AssertArity(mp, 2)
		return append(Associative{mp}, a...)
	}
	panic(ExpectedPair)
}

// Conjoin implements the Conjoiner interface
func (a Associative) Conjoin(v Value) Sequence {
	return a.Prepend(v)
}

// IsSequence returns whether this instance is a consumable Sequence
func (a Associative) IsSequence() bool {
	return len(a) > 0
}

// Str converts this Value into a Str
func (a Associative) Str() Str {
	var b bytes.Buffer
	l := len(a)

	b.WriteString("{")
	for i := 0; i < l; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		mp := a[i]
		b.WriteString(string(mp[0].Str()))
		b.WriteString(" ")
		b.WriteString(string(mp[1].Str()))
	}
	b.WriteString("}")
	return Str(b.String())
}

// AssertMapped will cast Value to a Mapped or explode violently
func AssertMapped(v Value) Mapped {
	if r, ok := v.(Mapped); ok {
		return r
	}
	panic(Err(ExpectedMapped, v))
}
