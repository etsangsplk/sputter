package builtins

import (
	"math/big"

	a "github.com/kode4food/sputter/api"
)

// NonNumericArgument is raised when a numeric argument is required
const NonNumericArgument = "arguments must be numeric"

type reduceFunc func(prev *big.Float, next *big.Float) *big.Float
type compareFunc func(prev *big.Float, next *big.Float) bool

func reduce(c a.Context, s a.Sequence, v *big.Float, f reduceFunc) a.Value {
	cur := v
	i := s.Iterate()
	for e, ok := i.Next(); ok; e, ok = i.Next() {
		if fv, ok := a.Eval(c, e).(*big.Float); ok {
			cur = f(cur, fv)
			continue
		}
		panic(NonNumericArgument)
	}
	return cur
}

func fetchFirstNumber(c a.Context, args a.Sequence) (*big.Float, a.Sequence) {
	a.AssertMinimumArity(args, 1)
	i := args.Iterate()
	v, _ := i.Next()
	if r, ok := a.Eval(c, v).(*big.Float); ok {
		return r, i.Rest()
	}
	panic(NonNumericArgument)
}

func compare(c a.Context, s a.Sequence, f compareFunc) a.Value {
	cur, r := fetchFirstNumber(c, s)
	i := r.Iterate()
	for e, ok := i.Next(); ok; e, ok = i.Next() {
		if v, ok := a.Eval(c, e).(*big.Float); ok {
			if !f(cur, v) {
				return a.False
			}
			cur = v
			continue
		}
		panic(NonNumericArgument)
	}
	return a.True
}

func add(c a.Context, args a.Sequence) a.Value {
	f := big.NewFloat(0)
	return reduce(c, args, f, func(p *big.Float, n *big.Float) *big.Float {
		return p.Add(p, n)
	})
}

func sub(c a.Context, args a.Sequence) a.Value {
	f, r := fetchFirstNumber(c, args)
	return reduce(c, r, f, func(p *big.Float, n *big.Float) *big.Float {
		return p.Sub(p, n)
	})
}

func mul(c a.Context, args a.Sequence) a.Value {
	f := big.NewFloat(1)
	return reduce(c, args, f, func(p *big.Float, n *big.Float) *big.Float {
		return p.Mul(p, n)
	})
}

func div(c a.Context, args a.Sequence) a.Value {
	f, r := fetchFirstNumber(c, args)
	return reduce(c, r, f, func(p *big.Float, n *big.Float) *big.Float {
		return p.Quo(p, n)
	})
}

func eq(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *big.Float, n *big.Float) bool {
		return p.Cmp(n) == 0
	})
}

func gt(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *big.Float, n *big.Float) bool {
		return p.Cmp(n) == 1
	})
}

func gte(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *big.Float, n *big.Float) bool {
		r := p.Cmp(n)
		return r == 0 || r == 1
	})
}

func lt(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *big.Float, n *big.Float) bool {
		return p.Cmp(n) == -1
	})
}

func lte(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *big.Float, n *big.Float) bool {
		r := p.Cmp(n)
		return r == 0 || r == -1
	})
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "+", Exec: add})
	a.PutFunction(Context, &a.Function{Name: "-", Exec: sub})
	a.PutFunction(Context, &a.Function{Name: "*", Exec: mul})
	a.PutFunction(Context, &a.Function{Name: "/", Exec: div})
	a.PutFunction(Context, &a.Function{Name: "=", Exec: eq})
	a.PutFunction(Context, &a.Function{Name: ">", Exec: gt})
	a.PutFunction(Context, &a.Function{Name: ">=", Exec: gte})
	a.PutFunction(Context, &a.Function{Name: "<", Exec: lt})
	a.PutFunction(Context, &a.Function{Name: "<=", Exec: lte})
}
