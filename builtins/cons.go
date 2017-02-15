package builtins

import a "github.com/kode4food/sputter/api"

// NonCons is thrown when you call car/cdr on a non-Cons
const NonCons = "Value is not a Cons cell"

func cons(c *a.Context, args a.Iterable) a.Value {
	AssertArity(args, 2)
	i := args.Iterate()
	car, _ := i.Next()
	cdr, _ := i.Next()
	return &a.Cons{Car: a.Evaluate(c, car), Cdr: a.Evaluate(c, cdr)}
}

func fetchCons(c *a.Context, args a.Iterable) *a.Cons {
	AssertArity(args, 1)
	i := args.Iterate()
	v, _ := i.Next()
	r := a.Evaluate(c, v)
	if cons, ok := r.(*a.Cons); ok {
		return cons
	}
	panic(NonCons)
}

func car(c *a.Context, args a.Iterable) a.Value {
	return fetchCons(c, args).Car
}

func cdr(c *a.Context, args a.Iterable) a.Value {
	return fetchCons(c, args).Cdr
}

func list(c *a.Context, args a.Iterable) a.Value {
	s := &a.Stack{}
	i := args.Iterate()
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		s.Push(a.Evaluate(c, v))
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

func isList(c *a.Context, args a.Iterable) a.Value {
	AssertArity(args, 1)
	i := args.Iterate()
	if v, ok := i.Next(); ok {
		if _, ok := a.Evaluate(c, v).(*a.Cons); ok {
			return a.True
		}
	}
	return a.False
}

func init() {
	Context.PutFunction(&a.Function{Name: "cons", Exec: cons})
	Context.PutFunction(&a.Function{Name: "car", Exec: car})
	Context.PutFunction(&a.Function{Name: "cdr", Exec: cdr})
	Context.PutFunction(&a.Function{Name: "list", Exec: list})
	Context.PutFunction(&a.Function{Name: "list?", Exec: isList})
}
