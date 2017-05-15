package api

import u "github.com/kode4food/sputter/util"

// ExpectedGetter is thrown if the Value is not a Getter
const ExpectedGetter = "expected a propertied value: %s"

var keywords = u.NewCache()

// Keyword is a Value that represents a Name that resolves to itself
type Keyword interface {
	Value
	Applicable
	Evaluable
	Named
	Keyword() bool
}

type keyword struct {
	name Name
}

// NewKeyword returns an interned instance of a Keyword
func NewKeyword(n Name) Keyword {
	return keywords.Get(n, func() u.Any {
		return &keyword{name: n}
	}).(Keyword)
}

// Keyword is a disambiguating marker
func (k *keyword) Keyword() bool {
	return true
}

// Name returns the Name component of the Keyword
func (k *keyword) Name() Name {
	return k.name
}

// Eval makes Keyword Evaluable
func (k *keyword) Eval(_ Context) Value {
	return k
}

// Apply makes Keyword Applicable
func (k *keyword) Apply(c Context, args Sequence) Value {
	AssertArity(args, 1)
	v := Eval(c, args.First())
	if g, ok := v.(Getter); ok {
		if r, ok := g.Get(k); ok {
			return r
		}
		panic(Err(KeyNotFound, k))
	}
	panic(Err(ExpectedGetter, v))
}

// Str converts this Value into a Str
func (k *keyword) Str() Str {
	return ":" + Str(k.name)
}
