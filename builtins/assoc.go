package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func assoc(_ a.Context, args a.Sequence) a.Value {
	return a.ToAssociative(args)
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
