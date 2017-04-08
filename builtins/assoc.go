package builtins

import a "github.com/kode4food/sputter/api"

func assoc(c a.Context, args a.Sequence) a.Value {
	if cnt, ok := args.(a.Countable); ok {
		l := cnt.Count()
		if l%2 != 0 {
			panic(a.ExpectedPair)
		}
		ml := l / 2
		r := make(a.Associative, ml)
		i := args
		for idx := 0; idx < ml; idx++ {
			k := a.Eval(c, i.First())
			i = i.Rest()

			v := a.Eval(c, i.First())
			i = i.Rest()

			r[idx] = a.Vector{k, v}
		}
		return r
	}
	return assocFromUncounted(c, args)
}

func assocFromUncounted(c a.Context, args a.Sequence) a.Value {
	r := a.Associative{}
	for i := args; i.IsSequence(); i = i.Rest() {
		k := i.First()
		i = i.Rest()
		if i.IsSequence() {
			v := i.First()
			r = append(r, a.Vector{
				a.Eval(c, k),
				a.Eval(c, v),
			})
		} else {
			panic(a.ExpectedPair)
		}
	}
	return r
}

func toAssoc(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	arg := a.Eval(c, args.First())
	seq := a.AssertSequence(arg)
	return assoc(c, seq)
}

func isAssociative(v a.Value) bool {
	if _, ok := v.(a.Associative); ok {
		return true
	}
	return false
}

func isMapped(v a.Value) bool {
	if _, ok := v.(a.Mapped); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(assoc).WithMetadata(a.Metadata{
			a.MetaName: a.Name("assoc"),
		}),
	)

	registerAnnotated(
		a.NewFunction(toAssoc).WithMetadata(a.Metadata{
			a.MetaName: a.Name("to-assoc"),
		}),
	)

	registerSequencePredicate(isAssociative, a.Metadata{
		a.MetaName: a.Name("assoc?"),
	})

	registerSequencePredicate(isMapped, a.Metadata{
		a.MetaName: a.Name("mapped?"),
	})
}
