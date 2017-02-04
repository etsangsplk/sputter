package sputter

import "math/big"
import "fmt"

// Builtins are the Context of built-in identifiers
var Builtins *Context

func addition(c *Context, args Iterable) Value {
	sum := big.NewFloat(0)
	next := args.Iterate()
	for value, found := next(); found; value, found = next() {
		sum.Add(sum, Evaluate(c, value).(*big.Float))
	}
	return sum
}

func subtraction(c *Context, args Iterable) Value {
	next := args.Iterate()
	value, found := next()
	diff := Evaluate(c, value).(*big.Float)
	for value, found = next(); found; value, found = next() {
		diff.Sub(diff, Evaluate(c, value).(*big.Float))
	}
	return diff
}

func multiplication(c *Context, args Iterable) Value {
	prod := big.NewFloat(1)
	next := args.Iterate()
	for value, found := next(); found; value, found = next() {
		prod.Mul(prod, Evaluate(c, value).(*big.Float))
	}
	return prod
}

func division(c *Context, args Iterable) Value {
	next := args.Iterate()
	value, found := next()
	quotient := Evaluate(c, value).(*big.Float)
	for value, found = next(); found; value, found = next() {
		quotient.Quo(quotient, Evaluate(c, value).(*big.Float))
	}
	return quotient
}

func puts(c *Context, args Iterable) Value {
	next := args.Iterate()
	for value, found := next(); found; {
		result := Evaluate(c, value)
		fmt.Print(result)
		value, found = next()
		if found {
			fmt.Print(" ")
		}
	}
	fmt.Println("")
	return nil
}

func init() {
	Builtins = NewContext()
	Builtins.PutFunction(&Function{"+", addition})
	Builtins.PutFunction(&Function{"-", subtraction})
	Builtins.PutFunction(&Function{"*", multiplication})
	Builtins.PutFunction(&Function{"/", division, })
	Builtins.PutFunction(&Function{"puts", puts})
}
