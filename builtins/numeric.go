package builtins

import (
	"math/big"

	a "github.com/kode4food/sputter/api"
)

func addition(c *a.Context, args a.Iterable) a.Value {
	sum := big.NewFloat(0)
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		sum.Add(sum, a.Evaluate(c, value).(*big.Float))
	}
	return sum
}

func subtraction(c *a.Context, args a.Iterable) a.Value {
	AssertMinimumArity(args, 1)
	i := args.Iterate()
	v, ok := i.Next()
	r := a.Evaluate(c, v).(*big.Float)
	for v, ok = i.Next(); ok; v, ok = i.Next() {
		r.Sub(r, a.Evaluate(c, v).(*big.Float))
	}
	return r
}

func multiplication(c *a.Context, args a.Iterable) a.Value {
	r := big.NewFloat(1)
	i := args.Iterate()
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		r.Mul(r, a.Evaluate(c, v).(*big.Float))
	}
	return r
}

func division(c *a.Context, args a.Iterable) a.Value {
	AssertMinimumArity(args, 1)
	i := args.Iterate()
	v, ok := i.Next()
	r := a.Evaluate(c, v).(*big.Float)
	for v, ok = i.Next(); ok; v, ok = i.Next() {
		r.Quo(r, a.Evaluate(c, v).(*big.Float))
	}
	return r
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "+", Exec: addition})
	BuiltIns.PutFunction(&a.Function{Name: "-", Exec: subtraction})
	BuiltIns.PutFunction(&a.Function{Name: "*", Exec: multiplication})
	BuiltIns.PutFunction(&a.Function{Name: "/", Exec: division})
}
