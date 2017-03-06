package builtins

import (
	a "github.com/kode4food/sputter/api"
	u "github.com/kode4food/sputter/util"
)

func cons(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	i := args.Iterate()
	car, _ := i.Next()
	cdr, _ := i.Next()
	return &a.Cons{Car: a.Eval(c, car), Cdr: a.Eval(c, cdr)}
}

func fetchCons(c a.Context, args a.Sequence) *a.Cons {
	a.AssertArity(args, 1)
	i := args.Iterate()
	v, _ := i.Next()
	return a.AssertCons(a.Eval(c, v))
}

func car(c a.Context, args a.Sequence) a.Value {
	return fetchCons(c, args).Car
}

func cdr(c a.Context, args a.Sequence) a.Value {
	return fetchCons(c, args).Cdr
}

func list(c a.Context, args a.Sequence) a.Value {
	s := u.NewStack()
	i := args.Iterate()
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		s.Push(a.Eval(c, v))
	}

	e, ok := s.Pop()
	if !ok {
		return a.Nil
	}

	var l = a.NewList(e)
	for v, ok := s.Pop(); ok; v, ok = s.Pop() {
		l = &a.Cons{Car: v, Cdr: l}
	}
	return l
}

func isList(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	i := args.Iterate()
	if v, ok := i.Next(); ok {
		if _, ok := a.Eval(c, v).(*a.Cons); ok {
			return a.True
		}
	}
	return a.False
}

func fetchSequence(c a.Context, args a.Sequence) a.Sequence {
	a.AssertArity(args, 1)
	i := args.Iterate()
	v, _ := i.Next()
	return a.AssertSequence(a.Eval(c, v))
}

func first(c a.Context, args a.Sequence) a.Value {
	s := fetchSequence(c, args)
	i := s.Iterate()
	v, _ := i.Next()
	return a.Eval(c, v)
}

func rest(c a.Context, args a.Sequence) a.Value {
	s := fetchSequence(c, args)
	i := s.Iterate()
	i.Next()
	return i.Rest()
}

func init() {
	putFunction(BuiltInNamespace, &a.Function{Name: "cons", Apply: cons})
	putFunction(BuiltInNamespace, &a.Function{Name: "car", Apply: car})
	putFunction(BuiltInNamespace, &a.Function{Name: "cdr", Apply: cdr})

	putFunction(BuiltInNamespace, &a.Function{Name: "list", Apply: list})
	registerPredicate(&a.Function{Name: "list?", Apply: isList})
	putFunction(BuiltInNamespace, &a.Function{Name: "first", Apply: first})
	putFunction(BuiltInNamespace, &a.Function{Name: "rest", Apply: rest})
}
