package sputter

import "fmt"

// Builtins are the Context of built-in identifiers
var Builtins *Context

func print(c *Context, args Iterable) Value {
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		result := Evaluate(c, value)
		fmt.Print(result)
	}
	fmt.Println("")
	return EmptyList
}

func defvar(c *Context, args Iterable) Value {
	globals := c.Globals()
	iter := args.Iterate()
	sym, _ := iter.Next()
	name := sym.(*Symbol).name
	_, bound := globals.Get(name)
	if !bound {
		val, _ := iter.Next()
		globals.Put(name, Evaluate(c, val))
	}
	return sym
}

func let(c *Context, args Iterable) Value {
	locals := c.Child()
	iter := args.Iterate()
	bindings, _ := iter.Next()

	bindIter := bindings.(Iterable).Iterate()
	for sym, ok := bindIter.Next(); ok; sym, ok = bindIter.Next() {
		name := sym.(*Symbol).name
		if val, ok := bindIter.Next(); ok {
			locals.Put(name, Evaluate(locals, val))
		}
	}
	
	return EvaluateIterator(locals, iter)
}

func defun(c *Context, args Iterable) Value {
	iter := args.Iterate()

	funcNameValue, _ := iter.Next()
	funcName := funcNameValue.(*Symbol).name

	symsValue, _ := iter.Next()
	syms := symsValue.(Iterable)

	body := iter.Iterable()

	defined := &Function{funcName, func(c *Context, args Iterable) Value {
		locals := c.Child()
		symIter := syms.Iterate()
		argIter := args.Iterate()
		for argSymbol, symFound := symIter.Next(); symFound; {
			argName := argSymbol.(*Symbol).name
			argValue, argFound := argIter.Next()
			if argFound {
				locals.Put(argName, argValue)
			} else {
				locals.Put(argName, EmptyList)
			}
			argSymbol, symFound = symIter.Next()
		}
		return EvaluateIterator(locals, body.Iterate())
	}}

	c.PutFunction(defined)
	return defined
}

func init() {
	Builtins = NewContext()
	Builtins.Put("T", &Literal{true})
	Builtins.Put("nil", EmptyList)
	Builtins.Put("true", &Literal{true})
	Builtins.Put("false", &Literal{false})

	Builtins.PutFunction(&Function{"print", print})
	Builtins.PutFunction(&Function{"defvar", defvar})
	Builtins.PutFunction(&Function{"let", let})
	Builtins.PutFunction(&Function{"defun", defun})
}
