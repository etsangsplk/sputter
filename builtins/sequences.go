package builtins

import a "github.com/kode4food/sputter/api"

func isSequence(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	if s, ok := a.Eval(c, v).(a.Sequence); ok {
		if s.IsSequence() {
			return a.True
		}
	}
	return a.False
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
	registerPredicate(
		a.NewFunction(isSequence).WithMetadata(a.Variables{
			a.MetaName: a.Name("seq?"),
		}).(a.Function),
	)

	registerAnnotated(
		a.NewFunction(first).WithMetadata(a.Variables{
			a.MetaName: a.Name("first"),
		}),
	)

	registerAnnotated(
		a.NewFunction(rest).WithMetadata(a.Variables{
			a.MetaName: a.Name("rest"),
		}),
	)

	registerAnnotated(
		a.NewFunction(cons).WithMetadata(a.Variables{
			a.MetaName: a.Name("cons"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_len).WithMetadata(a.Variables{
			a.MetaName: a.Name("len"),
		}),
	)

	registerAnnotated(
		a.NewFunction(concat).WithMetadata(a.Variables{
			a.MetaName: a.Name("concat"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_map).WithMetadata(a.Variables{
			a.MetaName: a.Name("map"),
		}),
	)

	registerAnnotated(
		a.NewFunction(filter).WithMetadata(a.Variables{
			a.MetaName: a.Name("filter"),
		}),
	)
}
