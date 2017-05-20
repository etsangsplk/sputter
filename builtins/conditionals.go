package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func _if(c a.Context, args a.Sequence) a.Value {
	i := a.AssertArityRange(args, 2, 3)
	cond := args.First().Eval(c)
	rest := args.Rest()
	if a.Truthy(cond) {
		return rest.First().Eval(c)
	}
	if i == 3 {
		return rest.Rest().First().Eval(c)
	}
	return a.Nil
}

func init() {
	registerAnnotated(
		a.NewFunction(_if).WithMetadata(a.Metadata{
			a.MetaName: a.Name("if"),
			a.MetaDoc:  d.Get("if"),
		}),
	)
}
