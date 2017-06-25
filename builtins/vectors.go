package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func vector(c a.Context, args a.Sequence) a.Value {
	if cnt, ok := args.(a.Counted); ok {
		l := cnt.Count()
		r := make([]a.Value, l)
		idx := 0
		for i := args; i.IsSequence(); i = i.Rest() {
			r[idx] = a.Eval(c, i.First())
			idx++
		}
		return a.NewVector(r...)
	}
	return vectorFromUncounted(c, args)
}

func vectorFromUncounted(c a.Context, args a.Sequence) a.Value {
	r := []a.Value{}
	for i := args; i.IsSequence(); i = i.Rest() {
		r = append(r, a.Eval(c, i.First()))
	}
	return a.NewVector(r...)
}

func isVector(v a.Value) bool {
	if _, ok := v.(a.Vector); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(vector).WithMetadata(a.Metadata{
			a.MetaName: a.Name("vector"),
			a.MetaDoc:  d.Get("vector"),
		}),
	)

	registerSequencePredicate(isVector, a.Metadata{
		a.MetaName: a.Name("vector?"),
		a.MetaDoc:  d.Get("is-vector"),
	})
}
