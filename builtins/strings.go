package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func str(_ a.Context, args a.Sequence) a.Value {
	return a.ToStr(args)
}

func isStr(v a.Value) bool {
	if _, ok := v.(a.Str); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(str).WithMetadata(a.Properties{
			a.MetaName: a.Name("str"),
			a.MetaDoc:  d.Get("str"),
		}),
	)

	registerSequencePredicate(isStr, a.Properties{
		a.MetaName: a.Name("str?"),
		a.MetaDoc:  d.Get("is-str"),
	})
}
