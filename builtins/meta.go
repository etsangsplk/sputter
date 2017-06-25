package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func toMetadata(args a.Mapped) a.Metadata {
	r := make(a.Metadata)
	for i := a.Sequence(args); i.IsSequence(); i = i.Rest() {
		p := i.First().(a.Sequence)
		k := p.First()
		v := p.Rest().First()
		r[k] = v
	}
	return r
}

func fromMetadata(m a.Metadata) a.Value {
	r := []a.Vector{}
	for k, v := range m {
		r = append(r, a.NewVector(k, v))
	}
	return a.NewAssociative(r...)
}

func withMeta(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	o := a.AssertAnnotated(args.First())
	m := a.AssertMapped(args.Rest().First())
	return o.WithMetadata(toMetadata(m)).(a.Value)
}

func meta(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	o := a.AssertAnnotated(args.First())
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
			a.MetaDoc:  d.Get("with-meta"),
		}),
	)

	registerAnnotated(
		a.NewFunction(meta).WithMetadata(a.Metadata{
			a.MetaName: a.Name("meta"),
			a.MetaDoc:  d.Get("meta"),
		}),
	)

	registerSequencePredicate(isAnnotated, a.Metadata{
		a.MetaName: a.Name("meta?"),
		a.MetaDoc:  d.Get("has-meta"),
	})
}
