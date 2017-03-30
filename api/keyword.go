package api

var keywords = make(Variables, 4096)

// Keyword is an Atom-like Value that represents a Name for mapping purposes
type Keyword interface {
	Applicable
	Evaluable
	Named
}

type keyword struct {
	name Name
}

// NewKeyword returns an interned instance of a Keyword
func NewKeyword(n Name) Keyword {
	if r, ok := keywords[n]; ok {
		return r.(Keyword)
	}
	r := &keyword{name: n}
	keywords[n] = r
	return r
}

// Name returns the Name component of the Keyword
func (k *keyword) Name() Name {
	return k.name
}

// Eval makes Keyword Evaluable
func (k *keyword) Eval(c Context) Value {
	return k
}

// Apply makes Keyword Applicable
func (k *keyword) Apply(c Context, args Sequence) Value {
	AssertArity(args, 1)
	if m, ok := Eval(c, args.First()).(Mapped); ok {
		return m.Get(k)
	}
	panic(ExpectedMapped)
}

func (k *keyword) String() string {
	return ":" + string(k.name)
}
