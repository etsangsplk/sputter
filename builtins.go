package sputter

import "math/big"
import "fmt"

// Builtins are the Context of built-in identifiers
var Builtins *Context

func addition(c *Context, args Iterable) Value {
	sum := big.NewFloat(0)
	iter := args.Iterate()
	for value, found := iter.Next(); found; value, found = iter.Next() {
		sum.Add(sum, Evaluate(c, value).(*big.Float))
	}
	return sum
}

func subtraction(c *Context, args Iterable) Value {
	iter := args.Iterate()
	value, found := iter.Next()
	diff := Evaluate(c, value).(*big.Float)
	for value, found = iter.Next(); found; value, found = iter.Next() {
		diff.Sub(diff, Evaluate(c, value).(*big.Float))
	}
	return diff
}

func multiplication(c *Context, args Iterable) Value {
	prod := big.NewFloat(1)
	iter := args.Iterate()
	for value, found := iter.Next(); found; value, found = iter.Next() {
		prod.Mul(prod, Evaluate(c, value).(*big.Float))
	}
	return prod
}

func division(c *Context, args Iterable) Value {
	iter := args.Iterate()
	value, found := iter.Next()
	quotient := Evaluate(c, value).(*big.Float)
	for value, found = iter.Next(); found; value, found = iter.Next() {
		quotient.Quo(quotient, Evaluate(c, value).(*big.Float))
	}
	return quotient
}

func print(c *Context, args Iterable) Value {
	iter := args.Iterate()
	for value, found := iter.Next(); found; {
		result := Evaluate(c, value)
		fmt.Print(result)
		value, found = iter.Next()
		if found {
			fmt.Print(" ")
		}
	}
	fmt.Println("")
	return EmptyList
}

func defvar(c *Context, args Iterable) Value {
	global := c.Global()
	iter := args.Iterate()
	sym, _ := iter.Next()
	name := sym.(*Symbol).name
	_, bound := global.Get(name)
	if !bound {
		val, _ := iter.Next()
		global.Put(name, Evaluate(c, val))
	}
	return sym
}

func bindSymbols(c *Context, iter Iterator) Value {
	var lastValue Value
	for sym, found := iter.Next(); found; sym, found = iter.Next() {
		if val, found := iter.Next(); found {
			lastValue = Evaluate(c, val)
			c.Put(sym.(*Symbol).name, lastValue)
		}
	}
	return lastValue
}

func let(c *Context, args Iterable) Value {
	locals := c.Child()
	iter := args.Iterate()
	bindings, _ := iter.Next()
	bindSymbols(locals, bindings.(Iterable).Iterate())
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

	Builtins.PutFunction(&Function{"+", addition})
	Builtins.PutFunction(&Function{"-", subtraction})
	Builtins.PutFunction(&Function{"*", multiplication})
	Builtins.PutFunction(&Function{"/", division})

	Builtins.PutFunction(&Function{"print", print})
	Builtins.PutFunction(&Function{"defvar", defvar})
	Builtins.PutFunction(&Function{"let", let})
	Builtins.PutFunction(&Function{"defun", defun})
}
