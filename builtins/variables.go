package builtins

import a "github.com/kode4food/sputter/api"

func defvarCommand(c *a.Context, args a.Iterable) a.Value {
	AssertMinimumArity(args, 2)
	globals := c.Globals()
	iter := args.Iterate()
	sym, _ := iter.Next()
	name := sym.(*a.Symbol).Name
	_, bound := globals.Get(name)
	if !bound {
		val, _ := iter.Next()
		globals.Put(name, a.Evaluate(c, val))
	}
	return sym
}

func letCommand(c *a.Context, args a.Iterable) a.Value {
	AssertMinimumArity(args, 2)
	locals := c.Child()
	iter := args.Iterate()
	bindings, _ := iter.Next()

	bindIter := bindings.(a.Iterable).Iterate()
	for sym, ok := bindIter.Next(); ok; sym, ok = bindIter.Next() {
		name := sym.(*a.Symbol).Name
		if val, ok := bindIter.Next(); ok {
			locals.Put(name, a.Evaluate(locals, val))
		}
	}

	return a.EvaluateIterator(locals, iter)
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "defvar", Exec: defvarCommand})
	BuiltIns.PutFunction(&a.Function{Name: "let", Exec: letCommand})
}
