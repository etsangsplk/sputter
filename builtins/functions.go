package builtins

import a "github.com/kode4food/sputter/api"

func define(name a.Name, argNames a.Iterable, body a.Iterable) *a.Function {
	ac := argCount(argNames)

	return &a.Function{
		Name: name,
		Exec: func(c *a.Context, args a.Iterable) a.Value {
			AssertArity(args, ac)
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
	g := c.Globals()
	i := args.Iterate()

	funcNameValue, _ := i.Next()
	funcName := funcNameValue.(*a.Symbol).Name

	argNamesValue, _ := i.Next()
	argNames := argNamesValue.(a.Iterable)

	body := i.Iterable()

	d := define(funcName, argNames, body)

	g.PutFunction(d)
	return d
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "defun", Exec: defunCommand})
}
