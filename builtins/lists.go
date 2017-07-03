package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func list(_ a.Context, args a.Sequence) a.Value {
	return a.ToList(args)
}

func isList(v a.Value) bool {
	if _, ok := v.(a.List); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(list).WithMetadata(a.Metadata{
			a.MetaName: a.Name("list"),
			a.MetaDoc:  d.Get("list"),
		}),
	)

	registerSequencePredicate(isList, a.Metadata{
		a.MetaName: a.Name("list?"),
		a.MetaDoc:  d.Get("is-list"),
	})
}
