package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func isSequence(v a.Value) bool {
	if s, ok := v.(a.Sequence); ok {
		return s.IsSequence()
	}
	return false
}

func fetchSequence(args a.Sequence) a.Sequence {
	a.AssertArity(args, 1)
	return a.AssertSequence(args.First())
}

func first(_ a.Context, args a.Sequence) a.Value {
	return fetchSequence(args).First()
}

func rest(_ a.Context, args a.Sequence) a.Value {
	return fetchSequence(args).Rest()
}

func cons(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f := args.First()
	r := args.Rest().First()
	return a.AssertSequence(r).Prepend(f)
}

func conj(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	s := a.AssertConjoiner(args.First())
	for i := args.Rest(); i.IsSequence(); i = i.Rest() {
		v := i.First()
		s = s.Conjoin(v).(a.Conjoiner)
	}
	return s
}

func _len(_ a.Context, args a.Sequence) a.Value {
	s := fetchSequence(args)
	l := a.Count(s)
	return a.NewFloat(float64(l))
}

func nth(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	s := a.AssertIndexed(args.First())
	return a.IndexedApply(s, c, args.Rest())
}

func _append(_ a.Context, args a.Sequence) a.Value {
	if a.AssertMinimumArity(args, 1) == 1 {
		r := args.First()
		return a.AssertSequence(r)
	}

	es := a.Map(args, func(v a.Value) a.Value {
		if v == a.Nil {
			return a.EmptyList
		}
		return a.AssertSequence(v)
	})

	return a.Concat(es)
}

func reduce(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f := a.AssertApplicable(args.First())
	s := _append(c, args.Rest()).(a.Sequence)
	return a.Reduce(c, s, f)
}

func take(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	n := a.AssertInteger(args.First())
	s := _append(c, args.Rest()).(a.Sequence)
	return a.Take(s, n)
}

func drop(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	n := a.AssertInteger(args.First())
	s := _append(c, args.Rest()).(a.Sequence)
	return a.Drop(s, n)
}

type forProc func(a.Context)

func forEach(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)

	b := a.AssertVector(args.First())
	bc := b.Count()
	if bc%2 != 0 {
		panic(ExpectedBindings)
	}

	var proc forProc
	depth := bc / 2
	for i := depth - 1; i >= 0; i-- {
		o := i * 2
		s, _ := b.ElementAt(o)
		e, _ := b.ElementAt(o + 1)
		n := a.AssertUnqualified(s).Name()
		if i == depth-1 {
			proc = makeTerminal(n, e, args.Rest())
		} else {
			proc = makeIntermediate(n, e, proc)
		}
	}

	proc(c)
	return a.Nil
}

func makeIntermediate(n a.Name, e a.Value, next forProc) forProc {
	return func(c a.Context) {
		s := a.AssertSequence(a.Eval(c, e))
		for i := s; i.IsSequence(); i = i.Rest() {
			l := a.ChildContext(c)
			l.Put(n, i.First())
			next(l)
		}
	}
}

func makeTerminal(n a.Name, e a.Value, bl a.Sequence) forProc {
	return func(c a.Context) {
		s := a.AssertSequence(a.Eval(c, e))
		for i := s; i.IsSequence(); i = i.Rest() {
			l := a.ChildContext(c)
			l.Put(n, i.First())
			a.EvalBlock(l, bl)
		}
	}
}

func init() {
	registerSequencePredicate(isSequence, a.Metadata{
		a.MetaName: a.Name("seq?"),
		a.MetaDoc:  d.Get("is-seq"),
	})

	registerAnnotated(
		a.NewFunction(first).WithMetadata(a.Metadata{
			a.MetaName: a.Name("first"),
			a.MetaDoc:  d.Get("first"),
		}),
	)

	registerAnnotated(
		a.NewFunction(rest).WithMetadata(a.Metadata{
			a.MetaName: a.Name("rest"),
			a.MetaDoc:  d.Get("rest"),
		}),
	)

	registerAnnotated(
		a.NewFunction(cons).WithMetadata(a.Metadata{
			a.MetaName: a.Name("cons"),
			a.MetaDoc:  d.Get("cons"),
		}),
	)

	registerAnnotated(
		a.NewFunction(conj).WithMetadata(a.Metadata{
			a.MetaName: a.Name("conj"),
			a.MetaDoc:  d.Get("conj"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_len).WithMetadata(a.Metadata{
			a.MetaName: a.Name("len"),
			a.MetaDoc:  d.Get("len"),
		}),
	)

	registerAnnotated(
		a.NewFunction(nth).WithMetadata(a.Metadata{
			a.MetaName: a.Name("nth"),
			a.MetaDoc:  d.Get("nth"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_append).WithMetadata(a.Metadata{
			a.MetaName: a.Name("append"),
		}),
	)

	registerAnnotated(
		a.NewFunction(reduce).WithMetadata(a.Metadata{
			a.MetaName: a.Name("reduce"),
			a.MetaDoc:  d.Get("reduce"),
		}),
	)

	registerAnnotated(
		a.NewFunction(take).WithMetadata(a.Metadata{
			a.MetaName: a.Name("take"),
			a.MetaDoc:  d.Get("take"),
		}),
	)

	registerAnnotated(
		a.NewFunction(drop).WithMetadata(a.Metadata{
			a.MetaName: a.Name("drop"),
			a.MetaDoc:  d.Get("drop"),
		}),
	)

	registerAnnotated(
		a.NewFunction(forEach).WithMetadata(a.Metadata{
			a.MetaName:    a.Name("for-each"),
			a.MetaSpecial: a.True,
		}),
	)
}
