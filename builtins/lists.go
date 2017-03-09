package builtins

import (
	a "github.com/kode4food/sputter/api"
	u "github.com/kode4food/sputter/util"
)

func cons(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f := a.Eval(c, args.First())
	r := a.Eval(c, args.Rest().First())
	return a.AssertSequence(r).Prepend(f)
}

func list(c a.Context, args a.Sequence) a.Value {
	s := u.NewStack()
	i := a.Iterate(args)
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		s.Push(a.Eval(c, v))
	}

	e, ok := s.Pop()
	if !ok {
		return a.EmptyList
	}

	var l = a.NewList(e)
	for v, ok := s.Pop(); ok; v, ok = s.Pop() {
		l = l.Prepend(v)
	}
	return l
}

func isList(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	if _, ok := a.Eval(c, v).(*a.List); ok {
		return a.True
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

func init() {
	registerFunction(&a.Function{Name: "cons", Apply: cons})
	registerFunction(&a.Function{Name: "list", Apply: list})
	registerPredicate(&a.Function{Name: "list?", Apply: isList})
	registerFunction(&a.Function{Name: "first", Apply: first})
	registerFunction(&a.Function{Name: "rest", Apply: rest})
}
