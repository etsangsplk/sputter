package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func assoc(_ a.Context, args a.Sequence) a.Value {
	if cnt, ok := args.(a.Counted); ok {
		l := cnt.Count()
		if l%2 != 0 {
			panic(a.ExpectedPair)
		}
		ml := l / 2
		r := make([]a.Vector, ml)
		i := args
		for idx := 0; idx < ml; idx++ {
			k := i.First()
			i = i.Rest()

			v := i.First()
			i = i.Rest()

			r[idx] = a.NewVector(k, v)
		}
		return a.NewAssociative(r...)
	}
	return assocFromUncounted(args)
}

func assocFromUncounted(args a.Sequence) a.Value {
	r := []a.Vector{}
	for i := args; i.IsSequence(); i = i.Rest() {
		k := i.First()
		i = i.Rest()
		if i.IsSequence() {
			v := i.First()
			r = append(r, a.NewVector(k, v))
		} else {
			panic(a.ExpectedPair)
		}
	}
	return a.NewAssociative(r...)
}

func isAssociative(v a.Value) bool {
	if _, ok := v.(a.Associative); ok {
		return true
	}
	return false
}

func isMapped(v a.Value) bool {
	if _, ok := v.(a.Mapped); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(assoc).WithMetadata(a.Metadata{
			a.MetaName: a.Name("assoc"),
			a.MetaDoc:  d.Get("assoc"),
		}),
	)

	registerSequencePredicate(isAssociative, a.Metadata{
		a.MetaName: a.Name("assoc?"),
		a.MetaDoc:  d.Get("is-assoc"),
	})

	registerSequencePredicate(isMapped, a.Metadata{
		a.MetaName: a.Name("mapped?"),
		a.MetaDoc:  d.Get("is-mapped"),
	})
}
