package builtins

import a "github.com/kode4food/sputter/api"

func isSequence(v a.Value) bool {
	if s, ok := v.(a.Sequence); ok {
		return s.IsSequence()
	}
	return false
}

func fetchSequence(c a.Context, args a.Sequence) a.Sequence {
	a.AssertArity(args, 1)
	return a.AssertSequence(a.Eval(c, args.First()))
}

func first(c a.Context, args a.Sequence) a.Value {
	return fetchSequence(c, args).First()
}

func rest(c a.Context, args a.Sequence) a.Value {
	return fetchSequence(c, args).Rest()
}

func cons(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f := a.Eval(c, args.First())
	r := a.Eval(c, args.Rest().First())
	return a.AssertSequence(r).Prepend(f)
}

func _len(c a.Context, args a.Sequence) a.Value {
	s := fetchSequence(c, args)
	l := a.Count(s)
	return a.NewFloat(float64(l))
}

func concat(c a.Context, args a.Sequence) a.Value {
	if a.AssertMinimumArity(args, 1) == 1 {
		return a.AssertSequence(a.Eval(c, args.First()))
	}

	es := a.NewMapper(args, func(v a.Value) a.Value {
		return a.AssertSequence(a.Eval(c, v))
	})

	return a.NewConcat(es)
}

func _map(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f := a.AssertApplicable(a.Eval(c, args.First()))
	s := concat(c, args.Rest()).(a.Sequence)
	return a.NewMapper(s, func(v a.Value) a.Value {
		return f.Apply(c, a.NewList(v))
	})
}

func filter(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f := a.AssertApplicable(a.Eval(c, args.First()))
	s := concat(c, args.Rest()).(a.Sequence)
	return a.NewFilter(s, func(v a.Value) bool {
		return a.Truthy(f.Apply(c, a.NewList(v)))
	})
}

func init() {
	registerSequencePredicate(isSequence, a.Metadata{
		a.MetaName: a.Name("seq?"),
	})

	registerAnnotated(
		a.NewFunction(first).WithMetadata(a.Metadata{
			a.MetaName: a.Name("first"),
		}),
	)

	registerAnnotated(
		a.NewFunction(rest).WithMetadata(a.Metadata{
			a.MetaName: a.Name("rest"),
		}),
	)

	registerAnnotated(
		a.NewFunction(cons).WithMetadata(a.Metadata{
			a.MetaName: a.Name("cons"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_len).WithMetadata(a.Metadata{
			a.MetaName: a.Name("len"),
		}),
	)

	registerAnnotated(
		a.NewFunction(concat).WithMetadata(a.Metadata{
			a.MetaName: a.Name("concat"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_map).WithMetadata(a.Metadata{
			a.MetaName: a.Name("map"),
		}),
	)

	registerAnnotated(
		a.NewFunction(filter).WithMetadata(a.Metadata{
			a.MetaName: a.Name("filter"),
		}),
	)
}
