package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

// BuiltIns is a special Context of built-in identifiers
var BuiltIns = a.NewContext()

func printCommand(c *a.Context, args a.Iterable) a.Value {
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		result := a.Evaluate(c, value)
		fmt.Print(result)
	}
	fmt.Println("")
	return a.EmptyList
}

func defvarCommand(c *a.Context, args a.Iterable) a.Value {
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

func defunCommand(c *a.Context, args a.Iterable) a.Value {
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

func ifCommand(c *a.Context, args a.Iterable) a.Value {
	iter := args.Iterate()
	condValue, _ := iter.Next()
	cond := a.Evaluate(c, condValue)
	if !a.Truthy(cond) {
		iter.Next()
	}
	result, _ := iter.Next()
	return a.Evaluate(c, result)
}

func isListCommand(c *a.Context, args a.Iterable) a.Value {
	iter := args.Iterate()
	if val, ok := iter.Next(); ok {
		if _, ok := a.Evaluate(c, val).(*a.List); ok {
			return a.True
		}
	}
	return a.False
}

func init() {
	BuiltIns.Put("T", a.True)
	BuiltIns.Put("nil", a.Nil)
	BuiltIns.Put("true", a.True)
	BuiltIns.Put("false", a.False)

	BuiltIns.PutFunction(&a.Function{Name: "print", Exec: printCommand})
	BuiltIns.PutFunction(&a.Function{Name: "defvar", Exec: defvarCommand})
	BuiltIns.PutFunction(&a.Function{Name: "let", Exec: letCommand})
	BuiltIns.PutFunction(&a.Function{Name: "defun", Exec: defunCommand})
	BuiltIns.PutFunction(&a.Function{Name: "if", Exec: ifCommand})
	BuiltIns.PutFunction(&a.Function{Name: "list?", Exec: isListCommand})
}
