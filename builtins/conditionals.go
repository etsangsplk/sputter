package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func cond(c a.Context, args a.Sequence) a.Value {
	for i := args; i.IsSequence(); i = i.Rest() {
		p := a.Eval(c, i.First())
		i = i.Rest()
		if !i.IsSequence() {
			return p
		}
		if a.Truthy(p) {
			return a.Eval(c, i.First())
		}
	}
	return a.Nil
}

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
		a.NewFunction(cond).WithMetadata(a.Metadata{
			a.MetaName: a.Name("cond"),
			a.MetaDoc:  d.Get("cond"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_if).WithMetadata(a.Metadata{
			a.MetaName: a.Name("if"),
			a.MetaDoc:  d.Get("if"),
		}),
	)
}
