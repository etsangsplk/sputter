package builtins

import a "github.com/kode4food/sputter/api"

func define(name string, argNames a.Iterable, body a.Iterable) *a.Function {
	count := argCount(argNames)

	return &a.Function{
		Name: name,
		Exec: func(c *a.Context, args a.Iterable) a.Value {
			AssertArity(args, count)
			locals := c.Child()
			argNamesIter := argNames.Iterate()
			argIter := args.Iterate()
			for nameSymbol, nameFound := argNamesIter.Next(); nameFound; {
				argName := nameSymbol.(*a.Symbol).Name
				argValue, argFound := argIter.Next()
				if argFound {
					locals.Put(argName, argValue)
				} else {
					locals.Put(argName, a.EmptyList)
				}
				nameSymbol, nameFound = argNamesIter.Next()
			}
			return a.EvaluateIterator(locals, body.Iterate())
		},
	}
}

func defunCommand(c *a.Context, args a.Iterable) a.Value {
	AssertMinimumArity(args, 3)
	globals := c.Globals()
	iter := args.Iterate()

	funcNameValue, _ := iter.Next()
	funcName := funcNameValue.(*a.Symbol).Name

	argNamesValue, _ := iter.Next()
	argNames := argNamesValue.(a.Iterable)

	body := iter.Iterable()

	defined := define(funcName, argNames, body)

	globals.PutFunction(defined)
	return defined
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "defun", Exec: defunCommand})
}
