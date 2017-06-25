package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func registerPredicate(f a.Function) {
	registerAnnotated(f)

	registerAnnotated(
		a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
			if f.Apply(c, args) == a.True {
				return a.False
			}
			return a.True
		}).WithMetadata(a.Metadata{
			a.MetaName: a.Name("!" + f.Name()),
		}),
	)
}

func registerSequencePredicate(f a.ValueFilter, md a.Metadata) {
	registerAnnotated(
		a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
			a.AssertMinimumArity(args, 1)
			for i := args; i.IsSequence(); i = i.Rest() {
				if !f(i.First()) {
					return a.False
				}
			}
			return a.True
		}).WithMetadata(md),
	)

	registerAnnotated(
		a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
			a.AssertMinimumArity(args, 1)
			for i := args; i.IsSequence(); i = i.Rest() {
				if f(i.First()) {
					return a.False
				}
			}
			return a.True
		}).WithMetadata(md.Merge(a.Metadata{
			a.MetaName: a.Name("!" + md[a.MetaName].(a.Name)),
		})),
	)
}

func identical(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	l := args.First()
	for i := args.Rest(); i.IsSequence(); i = i.Rest() {
		if l != i.First() {
			return a.False
		}
	}
	return a.True
}

func init() {
	registerPredicate(
		a.NewFunction(identical).WithMetadata(a.Metadata{
			a.MetaName: a.Name("eq"),
			a.MetaDoc:  d.Get("eq"),
		}).(a.Function),
	)

	registerSequencePredicate(func(v a.Value) bool {
		return v == a.Nil
	}, a.Metadata{
		a.MetaName: a.Name("nil?"),
		a.MetaDoc:  d.Get("is-nil"),
	})
}
