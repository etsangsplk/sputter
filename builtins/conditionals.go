package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func _if(c a.Context, args a.Sequence) a.Value {
	i := a.AssertArityRange(args, 2, 3)
	cond := a.Eval(c, args.First())
	rest := args.Rest()
	if a.Truthy(cond) {
		return a.Eval(c, rest.First())
	}
	if i == 3 {
		return a.Eval(c, rest.Rest().First())
	}
	return a.Nil
}

func init() {
	registerAnnotated(
		a.NewFunction(_if).WithMetadata(a.Properties{
			a.MetaName:    a.Name("if"),
			a.MetaDoc:     d.Get("if"),
			a.MetaSpecial: a.True,
		}),
	)
}
