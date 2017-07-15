package builtins

import (
	a "github.com/kode4food/sputter/api"
)

func makeLazySequence(c a.Context, args a.Sequence) a.Value {
	db := a.NewBlock(args)

	f := a.NewFunction(func(c a.Context, _ a.Sequence) a.Value {
		return a.Eval(c, db)
	})

	return a.NewLazySequence(a.MakeLazyResolver(c, f))
}

func concat(_ a.Context, args a.Sequence) a.Value {
	if a.AssertMinimumArity(args, 1) == 1 {
		return a.AssertSequence(args.First())
	}
	return a.Concat(args)
}

func filter(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f := a.AssertApplicable(args.First())
	s := a.Concat(args.Rest())
	return a.Filter(s, a.MakeValueFilter(c, f))
}

func _map(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f := a.AssertApplicable(args.First())
	s := a.Concat(args.Rest())
	return a.Map(s, a.MakeValueMapper(c, f))
}

func reduce(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f := a.AssertApplicable(args.First())
	s := a.Concat(args.Rest())
	return a.Reduce(s, a.MakeValueReducer(c, f))
}

func take(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	n := a.AssertInteger(args.First())
	s := a.Concat(args.Rest())
	return a.Take(s, n)
}

func drop(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	n := a.AssertInteger(args.First())
	s := a.Concat(args.Rest())
	return a.Drop(s, n)
}

func makeRange(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 3)
	low := a.AssertNumber(args.First())
	high := a.AssertNumber(args.Rest().First())
	step := a.AssertNumber(args.Rest().Rest().First())
	return a.NewRange(low, high, step)
}

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
	RegisterBuiltIn("make-lazy-seq", makeLazySequence)
	RegisterBuiltIn("concat", concat)
	RegisterBuiltIn("filter", filter)
	RegisterBuiltIn("map", _map)
	RegisterBuiltIn("reduce", reduce)
	RegisterBuiltIn("take", take)
	RegisterBuiltIn("drop", drop)
	RegisterBuiltIn("make-range", makeRange)
	RegisterBuiltIn("for-each", forEach)
}
