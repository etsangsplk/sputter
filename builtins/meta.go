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

func fromMetadata(m a.Metadata) a.Value {
	r := a.ArrayMap{}
	for k, v := range m {
		r = r.Prepend(a.Vector{k, v}).(a.ArrayMap)
	}
	return r
}

func withMeta(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	o := a.AssertAnnotated(a.Eval(c, args.First()))
	m := a.AssertSequence(a.Eval(c, args.Rest().First()))
	return o.WithMetadata(toMetadata(m))
}

func meta(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	o := a.AssertAnnotated(a.Eval(c, args.First()))
	return fromMetadata(o.Metadata())
}

func isAnnotated(v a.Value) bool {
	if _, ok := v.(a.Annotated); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(withMeta).WithMetadata(a.Metadata{
			a.MetaName: a.Name("with-meta"),
		}),
	)

	registerAnnotated(
		a.NewFunction(meta).WithMetadata(a.Metadata{
			a.MetaName: a.Name("meta"),
		}),
	)

	registerSequencePredicate(isAnnotated, a.Metadata{
		a.MetaName: a.Name("meta?"),
	})
}
