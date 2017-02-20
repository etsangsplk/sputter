package builtins

import a "github.com/kode4food/sputter/api"

// NonCons is thrown when you call car/cdr on a non-Cons
const NonCons = "Value is not a Cons cell"

func cons(c a.Context, args a.Sequence) a.Value {
	AssertArity(args, 2)
	i := args.Iterate()
	car, _ := i.Next()
	cdr, _ := i.Next()
	return &a.Cons{Car: a.Eval(c, car), Cdr: a.Eval(c, cdr)}
}

func fetchCons(c a.Context, args a.Sequence) *a.Cons {
	AssertArity(args, 1)
	i := args.Iterate()
	v, _ := i.Next()
	r := a.Eval(c, v)
	if cons, ok := r.(*a.Cons); ok {
		return cons
	}
	panic(NonCons)
}

func car(c a.Context, args a.Sequence) a.Value {
	return fetchCons(c, args).Car
}

func cdr(c a.Context, args a.Sequence) a.Value {
	return fetchCons(c, args).Cdr
}

func list(c a.Context, args a.Sequence) a.Value {
	s := a.NewStack()
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
	AssertArity(args, 1)
	i := args.Iterate()
	if v, ok := i.Next(); ok {
		if _, ok := a.Eval(c, v).(*a.Cons); ok {
			return a.True
		}
	}
	return a.False
}

func first(c a.Context, args a.Sequence) a.Value {
	AssertArity(args, 1)
	return fetchCons(c, args).Get(0)
}

func second(c a.Context, args a.Sequence) a.Value {
	AssertArity(args, 1)
	return fetchCons(c, args).Get(1)
}

func third(c a.Context, args a.Sequence) a.Value {
	AssertArity(args, 1)
	return fetchCons(c, args).Get(2)
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "cons", Exec: cons})
	a.PutFunction(Context, &a.Function{Name: "car", Exec: car})
	a.PutFunction(Context, &a.Function{Name: "cdr", Exec: cdr})
	a.PutFunction(Context, &a.Function{Name: "list", Exec: list})
	a.PutFunction(Context, &a.Function{Name: "list?", Exec: isList})
	a.PutFunction(Context, &a.Function{Name: "first", Exec: first})
	a.PutFunction(Context, &a.Function{Name: "second", Exec: second})
	a.PutFunction(Context, &a.Function{Name: "third", Exec: third})
}
