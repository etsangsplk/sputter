package builtins

import (
	a "github.com/kode4food/sputter/api"
)

func toMetadata(args a.Sequence) a.Metadata {
	r := make(a.Metadata)
	for i := args.(a.Sequence); i.IsSequence(); i = i.Rest() {
		p := a.AssertSequence(i.First())
		a.AssertArity(p, 2)
		k := p.First()
		v := p.Rest().First()
		r[k] = v
	}
	return r
}

func withMeta(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	o := a.AssertAnnotated(a.Eval(c, args.First()))
	m := a.AssertSequence(a.Eval(c, args.Rest().First()))
	return o.WithMetadata(toMetadata(m))
}

func init() {
	registerAnnotated(
		a.NewFunction(withMeta).WithMetadata(a.Metadata{
			a.MetaName: a.Name("with-meta"),
		}),
	)
}
