package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func list(c a.Context, args a.Sequence) a.Value {
	if cnt, ok := args.(a.Counted); ok {
		l := cnt.Count()
		r := make([]a.Value, l)
		idx := 0
		for i := args; i.IsSequence(); i = i.Rest() {
			r[idx] = a.Eval(c, i.First())
			idx++
		}
		return a.NewList(r...)
	}
	return listFromUncounted(c, args)
}

func listFromUncounted(c a.Context, args a.Sequence) a.Value {
	r := []a.Value{}
	for i := args; i.IsSequence(); i = i.Rest() {
		r = append(r, a.Eval(c, i.First()))
	}
	return a.NewList(r...)
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
