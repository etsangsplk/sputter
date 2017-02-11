package builtins

import a "github.com/kode4food/sputter/api"

func defunCommand(c *a.Context, args a.Iterable) a.Value {
	AssertMinimumArity(args, 3)
	globals := c.Globals()
	iter := args.Iterate()

	funcNameValue, _ := iter.Next()
	funcName := funcNameValue.(*a.Symbol).Name

	symsValue, _ := iter.Next()
	syms := symsValue.(a.Iterable)

	body := iter.Iterable()

	defined := &a.Function{
		Name: funcName,
		Exec: func(c *a.Context, args a.Iterable) a.Value {
			locals := c.Child()
			symIter := syms.Iterate()
			argIter := args.Iterate()
			for argSymbol, symFound := symIter.Next(); symFound; {
				argName := argSymbol.(*a.Symbol).Name
				argValue, argFound := argIter.Next()
				if argFound {
					locals.Put(argName, argValue)
				} else {
					locals.Put(argName, a.EmptyList)
				}
				argSymbol, symFound = symIter.Next()
			}
			return a.EvaluateIterator(locals, body.Iterate())
		},
	}

	globals.PutFunction(defined)
	return defined
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "defun", Exec: defunCommand})
}
