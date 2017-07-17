package api

import "bytes"

const (
	// ExpectedPair is thrown if you prepend to a Map incorrectly
	ExpectedPair = "expected two element vectors when prepending"

	// ExpectedMapped is thrown if the Value is not a Mapped item
	ExpectedMapped = "expected a mapped value: %s"

	// KeyNotFound is thrown if a Key is not found in a Mapped item
	KeyNotFound = "key not found: %s"
)

type (
	// Associative is a Mappable that is implemented atop an array
	Associative interface {
		Conjoiner
		Mapped
		Counted
		Applicable
		Evaluable
		IsAssociative() bool
	}

	associative []Vector
)

// NewAssociative instantiates a new Associative
func NewAssociative(v ...Vector) Associative {
	return associative(v)
}

func (a associative) IsAssociative() bool {
	return true
}

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

func (a associative) Apply(_ Context, args Sequence) Value {
	AssertArity(args, 1)
	k := args.First()
	if r, ok := a.Get(k); ok {
		return r
	}
	panic(Err(KeyNotFound, k))
}

func (a associative) Eval(c Context) Value {
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

func (a associative) First() Value {
	return a[0]
}

func (a associative) Rest() Sequence {
	return Sequence(a[1:])
}

func (a associative) Prepend(v Value) Sequence {
	if mp, ok := v.(Vector); ok {
		AssertArity(mp, 2)
		return append(associative{mp}, a...)
	}
	panic(ExpectedPair)
}

func (a associative) Conjoin(v Value) Sequence {
	return a.Prepend(v)
}

func (a associative) IsSequence() bool {
	return len(a) > 0
}

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

// MappedApply provides 'get' behavior for Mapped Values
func MappedApply(s Mapped, args Sequence) Value {
	i := AssertArityRange(args, 1, 2)
	key := args.First()
	if r, ok := s.Get(key); ok {
		return r
	}
	if i == 2 {
		return args.Rest().First()
	}
	panic(Err(KeyNotFound, key))
}

// AssertMapped will cast Value to a Mapped or explode violently
func AssertMapped(v Value) Mapped {
	if r, ok := v.(Mapped); ok {
		return r
	}
	panic(Err(ExpectedMapped, v))
}
