package api

import u "github.com/kode4food/sputter/util"

var keywords = u.NewCache()

// Keyword is a Value that represents a Name that resolves to itself
type Keyword interface {
	Value
	Applicable
	Name() Name
	KeywordType()
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

func (k *keyword) KeywordType() {}

func (k *keyword) Name() Name {
	return k.name
}

func (k *keyword) Apply(_ Context, args Sequence) Value {
	i := AssertArityRange(args, 1, 2)
	s := args.First().(Mapped)
	if r, ok := s.Get(k); ok {
		return r
	}
	if i == 2 {
		return args.Rest().First()
	}
	panic(ErrStr(KeyNotFound, k))
}

func (k *keyword) Str() Str {
	return ":" + Str(k.name)
}
