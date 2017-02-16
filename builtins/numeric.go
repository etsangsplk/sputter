package builtins

import (
	"math/big"

	a "github.com/kode4food/sputter/api"
)

func add(c *a.Context, args a.Sequence) a.Value {
	sum := big.NewFloat(0)
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		sum.Add(sum, a.Eval(c, value).(*big.Float))
	}
	return sum
}

func sub(c *a.Context, args a.Sequence) a.Value {
	AssertMinimumArity(args, 1)
	i := args.Iterate()
	v, ok := i.Next()
	r := a.Eval(c, v).(*big.Float)
	for v, ok = i.Next(); ok; v, ok = i.Next() {
		r.Sub(r, a.Eval(c, v).(*big.Float))
	}
	return r
}

func mul(c *a.Context, args a.Sequence) a.Value {
	r := big.NewFloat(1)
	i := args.Iterate()
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		r.Mul(r, a.Eval(c, v).(*big.Float))
	}
	return r
}

func div(c *a.Context, args a.Sequence) a.Value {
	AssertMinimumArity(args, 1)
	i := args.Iterate()
	v, ok := i.Next()
	r := a.Eval(c, v).(*big.Float)
	for v, ok = i.Next(); ok; v, ok = i.Next() {
		r.Quo(r, a.Eval(c, v).(*big.Float))
	}
	return r
}

func init() {
	Context.PutFunction(&a.Function{Name: "+", Exec: add})
	Context.PutFunction(&a.Function{Name: "-", Exec: sub})
	Context.PutFunction(&a.Function{Name: "*", Exec: mul})
	Context.PutFunction(&a.Function{Name: "/", Exec: div})
}
