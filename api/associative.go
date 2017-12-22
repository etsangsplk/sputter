package api

import "bytes"

const (
	// ExpectedPair is thrown if you prepend to a Map incorrectly
	ExpectedPair = "expected a key-value pair"

	// KeyNotFound is thrown if a Key is not found in a Mapped item
	KeyNotFound = "key not found: %s"
)

type (
	// Associative is a Mappable that is implemented atop an Values
	Associative interface {
		Conjoiner
		Mapped
		Counted
		Applicable
		Evaluable
		AssociativeType()
	}

	associative []Vector
)

var emptyAssociative = associative{}

// NewAssociative instantiates a new Associative
func NewAssociative(v ...Vector) Associative {
	return associative(v)
}

func (associative) AssociativeType() {}

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
	return MappedApply(a, args)
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
	if len(a) > 0 {
		return a[0]
	}
	return Nil
}

func (a associative) Rest() Sequence {
	if len(a) > 1 {
		return a[1:]
	}
	return emptyAssociative
}

func (a associative) Split() (Value, Sequence, bool) {
	if len(a) > 0 {
		return a[0], a[1:], true
	}
	return Nil, emptyAssociative, false
}

func (a associative) Prepend(v Value) Sequence {
	if mp, ok := v.(Vector); ok && mp.Count() == 2 {
		return append(associative{mp}, a...)
	}
	panic(ErrStr(ExpectedPair))
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
	panic(ErrStr(KeyNotFound, key))
}
