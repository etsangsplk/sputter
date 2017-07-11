package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func vector(_ a.Context, args a.Sequence) a.Value {
	return a.ToVector(args)
}

func isVector(v a.Value) bool {
	if _, ok := v.(a.Vector); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(vector).WithMetadata(a.Properties{
			a.MetaName: a.Name("vector"),
			a.MetaDoc:  d.Get("vector"),
		}),
	)

	registerSequencePredicate(isVector, a.Properties{
		a.MetaName: a.Name("vector?"),
		a.MetaDoc:  d.Get("is-vector"),
	})
}
