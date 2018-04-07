package api

import "bytes"

const (
	// ExpectedPair is thrown if you prepend to a Map incorrectly
	ExpectedPair = "expected a key-value pair"

	// KeyNotFound is thrown if a Key is not found in a Mapped item
	KeyNotFound = "key not found: %s"
)

// Associative is a Mapped Value that is implemented atop an Vector
type Associative []Vector

// EmptyAssociative represents an empty Associative
var EmptyAssociative = Associative{}

// NewAssociative instantiates a new Associative
func NewAssociative(v ...Vector) Associative {
	return Associative(v)
}

// Count returns the number of pairs in the Associative
func (a Associative) Count() int {
	return len(a)
}

// Get returns the Value corresponding to the key in the Associative
func (a Associative) Get(key Value) (Value, bool) {
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

// First returns the first pair of the Associative
func (a Associative) First() Value {
	if len(a) > 0 {
		return a[0]
	}
	return Nil
}

// Rest returns the pairs of the List that follow the first
func (a Associative) Rest() Sequence {
	if len(a) > 1 {
		return a[1:]
	}
	return EmptyAssociative
}

// IsSequence returns whether or not this Associative has any pairs
func (a Associative) IsSequence() bool {
	return len(a) > 0
}

// Split breaks the Associative into its components (first, rest, isSequence)
func (a Associative) Split() (Value, Sequence, bool) {
	if len(a) > 0 {
		return a[0], a[1:], true
	}
	return Nil, EmptyAssociative, false
}

// Prepend inserts a pair at the beginning of the Associative
func (a Associative) Prepend(v Value) Sequence {
	if mp, ok := v.(Vector); ok && mp.Count() == 2 {
		return append(Associative{mp}, a...)
	}
	panic(ErrStr(ExpectedPair))
}

// Conjoin inserts a pair at the beginning of the Associative
func (a Associative) Conjoin(v Value) Sequence {
	return a.Prepend(v)
}

// Apply makes Associative Applicable
func (a Associative) Apply(_ Context, args Vector) Value {
	return MappedApply(a, args)
}

// Eval evaluates its elements, returning a new Associative
func (a Associative) Eval(c Context) Value {
	l := len(a)
	r := make(Associative, l)
	for i := 0; i < l; i++ {
		mp := a[i]
		k, _ := mp.ElementAt(0)
		v, _ := mp.ElementAt(1)
		r[i] = NewVector(Eval(c, k), Eval(c, v))
	}
	return r
}

// Str converts this Associative to a Str
func (a Associative) Str() Str {
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

// MappedApply provides 'get' behavior for Mapped Vector
func MappedApply(s Mapped, args Vector) Value {
	i := AssertArityRange(args, 1, 2)
	key := args[0]
	if r, ok := s.Get(key); ok {
		return r
	}
	if i == 2 {
		return args[1]
	}
	panic(ErrStr(KeyNotFound, key))
}
