package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
	i "github.com/kode4food/sputter/interpreter"
)

// BuiltIns is a special Context of built-in identifiers
var BuiltIns = a.NewContext()

func print(c *a.Context, args a.Iterable) a.Value {
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		result := i.Evaluate(c, value)
		fmt.Print(result)
	}
	fmt.Println("")
	return a.EmptyList
}

func defvar(c *a.Context, args a.Iterable) a.Value {
	globals := c.Globals()
	iter := args.Iterate()
	sym, _ := iter.Next()
	name := sym.(*a.Symbol).Name
	_, bound := globals.Get(name)
	if !bound {
		val, _ := iter.Next()
		globals.Put(name, i.Evaluate(c, val))
	}
	return sym
}

func let(c *a.Context, args a.Iterable) a.Value {
	locals := c.Child()
	iter := args.Iterate()
	bindings, _ := iter.Next()

	bindIter := bindings.(a.Iterable).Iterate()
	for sym, ok := bindIter.Next(); ok; sym, ok = bindIter.Next() {
		name := sym.(*a.Symbol).Name
		if val, ok := bindIter.Next(); ok {
			locals.Put(name, i.Evaluate(locals, val))
		}
	}

	return i.EvaluateIterator(locals, iter)
}

func defun(c *a.Context, args a.Iterable) a.Value {
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
			return i.EvaluateIterator(locals, body.Iterate())
		},
	}

	globals.PutFunction(defined)
	return defined
}

func init() {
	BuiltIns.Put("T", &a.Literal{Value: true})
	BuiltIns.Put("nil", a.EmptyList)
	BuiltIns.Put("true", &a.Literal{Value: true})
	BuiltIns.Put("false", &a.Literal{Value: false})

	BuiltIns.PutFunction(&a.Function{Name: "print", Exec: print})
	BuiltIns.PutFunction(&a.Function{Name: "defvar", Exec: defvar})
	BuiltIns.PutFunction(&a.Function{Name: "let", Exec: let})
	BuiltIns.PutFunction(&a.Function{Name: "defun", Exec: defun})
}
