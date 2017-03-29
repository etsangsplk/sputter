package api

var keywords = make(Variables, 4096)

// Keyword is an Atom-like Value that represents a Name for mapping purposes
type Keyword struct {
	name Name
}

// NewKeyword returns an interned instance of a Keyword
func NewKeyword(n Name) *Keyword {
	if r, ok := keywords[n]; ok {
		return r.(*Keyword)
	}
	r := &Keyword{name: n}
	keywords[n] = r
	return r
}

func (k *Keyword) Name() Name {
	return k.name
}

// Eval makes Keyword Evaluable
func (k *Keyword) Eval(c Context) Value {
	return k
}

// Apply makes Keyword Applicable
func (k *Keyword) Apply(c Context, args Sequence) Value {
	AssertArity(args, 1)
	if m, ok := Eval(c, args.First()).(Mapped); ok {
		return m.Get(k)
	}
	panic(ExpectedMapped)
}

func (k *Keyword) String() string {
	return ":" + string(k.name)
}
