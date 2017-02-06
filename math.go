package sputter

import "math/big"

func addition(c *Context, args Iterable) Value {
	sum := big.NewFloat(0)
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		sum.Add(sum, Evaluate(c, value).(*big.Float))
	}
	return sum
}

func subtraction(c *Context, args Iterable) Value {
	iter := args.Iterate()
	value, ok := iter.Next()
	diff := Evaluate(c, value).(*big.Float)
	for value, ok = iter.Next(); ok; value, ok = iter.Next() {
		diff.Sub(diff, Evaluate(c, value).(*big.Float))
	}
	return diff
}

func multiplication(c *Context, args Iterable) Value {
	prod := big.NewFloat(1)
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		prod.Mul(prod, Evaluate(c, value).(*big.Float))
	}
	return prod
}

func division(c *Context, args Iterable) Value {
	iter := args.Iterate()
	value, ok := iter.Next()
	quotient := Evaluate(c, value).(*big.Float)
	for value, ok = iter.Next(); ok; value, ok = iter.Next() {
		quotient.Quo(quotient, Evaluate(c, value).(*big.Float))
	}
	return quotient
}

func init() {
	Builtins.PutFunction(&Function{"+", addition})
	Builtins.PutFunction(&Function{"-", subtraction})
	Builtins.PutFunction(&Function{"*", multiplication})
	Builtins.PutFunction(&Function{"/", division})
}
