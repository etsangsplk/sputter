package builtins

import a "github.com/kode4food/sputter/api"

func cons(c *a.Context, args a.Iterable) a.Value {
	AssertArity(args, 2)
	i := args.Iterate()
	car, _ := i.Next()
	cdr, _ := i.Next()
	return &a.Cons{Car: a.Evaluate(c, car), Cdr: a.Evaluate(c, cdr)}
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
	Context.PutFunction(&a.Function{Name: "list?", Exec: isList})
}
