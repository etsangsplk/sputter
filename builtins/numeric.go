package builtins

import (
	"math/big"

	a "github.com/kode4food/sputter/api"
	i "github.com/kode4food/sputter/interpreter"
)

func addition(c *a.Context, args a.Iterable) a.Value {
	sum := big.NewFloat(0)
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		sum.Add(sum, i.Evaluate(c, value).(*big.Float))
	}
	return sum
}

func subtraction(c *a.Context, args a.Iterable) a.Value {
	iter := args.Iterate()
	value, ok := iter.Next()
	diff := i.Evaluate(c, value).(*big.Float)
	for value, ok = iter.Next(); ok; value, ok = iter.Next() {
		diff.Sub(diff, i.Evaluate(c, value).(*big.Float))
	}
	return diff
}

func multiplication(c *a.Context, args a.Iterable) a.Value {
	prod := big.NewFloat(1)
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		prod.Mul(prod, i.Evaluate(c, value).(*big.Float))
	}
	return prod
}

func division(c *a.Context, args a.Iterable) a.Value {
	iter := args.Iterate()
	value, ok := iter.Next()
	quotient := i.Evaluate(c, value).(*big.Float)
	for value, ok = iter.Next(); ok; value, ok = iter.Next() {
		quotient.Quo(quotient, i.Evaluate(c, value).(*big.Float))
	}
	return quotient
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "+", Exec: addition})
	BuiltIns.PutFunction(&a.Function{Name: "-", Exec: subtraction})
	BuiltIns.PutFunction(&a.Function{Name: "*", Exec: multiplication})
	BuiltIns.PutFunction(&a.Function{Name: "/", Exec: division})
}
