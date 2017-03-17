package builtins

import a "github.com/kode4food/sputter/api"

type reduceFunc func(prev *a.Number, next *a.Number) *a.Number
type compareFunc func(prev *a.Number, next *a.Number) bool

var (
	zero = a.NewFloat(0)
	one  = a.NewFloat(1)
)

func reduce(c a.Context, s a.Sequence, v *a.Number, f reduceFunc) a.Value {
	cur := v
	i := a.Iterate(s)
	for e, ok := i.Next(); ok; e, ok = i.Next() {
		fv := a.AssertNumber(a.Eval(c, e))
		cur = f(cur, fv)
	}
	return cur
}

func fetchFirstNumber(c a.Context, args a.Sequence) (*a.Number, a.Sequence) {
	a.AssertMinimumArity(args, 1)
	i := a.Iterate(args)
	v, _ := i.Next()
	nv := a.AssertNumber(a.Eval(c, v))
	return nv, i.Rest()
}

func compare(c a.Context, s a.Sequence, f compareFunc) a.Value {
	cur, r := fetchFirstNumber(c, s)
	i := a.Iterate(r)
	for e, ok := i.Next(); ok; e, ok = i.Next() {
		v := a.AssertNumber(a.Eval(c, e))
		if !f(cur, v) {
			return a.False
		}
		cur = v
	}
	return a.True
}

func add(c a.Context, args a.Sequence) a.Value {
	if !args.IsSequence() {
		return zero
	}
	f, r := fetchFirstNumber(c, args)
	return reduce(c, r, f, func(p *a.Number, n *a.Number) *a.Number {
		return p.Add(n)
	})
}

func sub(c a.Context, args a.Sequence) a.Value {
	f, r := fetchFirstNumber(c, args)
	return reduce(c, r, f, func(p *a.Number, n *a.Number) *a.Number {
		return p.Sub(n)
	})
}

func mul(c a.Context, args a.Sequence) a.Value {
	if !args.IsSequence() {
		return one
	}
	f, r := fetchFirstNumber(c, args)
	return reduce(c, r, f, func(p *a.Number, n *a.Number) *a.Number {
		return p.Mul(n)
	})
}

func div(c a.Context, args a.Sequence) a.Value {
	f, r := fetchFirstNumber(c, args)
	return reduce(c, r, f, func(p *a.Number, n *a.Number) *a.Number {
		return p.Div(n)
	})
}

func eq(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *a.Number, n *a.Number) bool {
		return p.Cmp(n) == a.EqualTo
	})
}

func neq(c a.Context, args a.Sequence) a.Value {
	if eq(c, args) == a.True {
		return a.False
	}
	return a.True
}

func gt(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *a.Number, n *a.Number) bool {
		return p.Cmp(n) == a.GreaterThan
	})
}

func gte(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *a.Number, n *a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.GreaterThan
	})
}

func lt(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *a.Number, n *a.Number) bool {
		return p.Cmp(n) == a.LessThan
	})
}

func lte(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p *a.Number, n *a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.LessThan
	})
}

func init() {
	registerFunction(&a.Function{Name: "+", Exec: add})
	registerFunction(&a.Function{Name: "-", Exec: sub})
	registerFunction(&a.Function{Name: "*", Exec: mul})
	registerFunction(&a.Function{Name: "/", Exec: div})
	registerFunction(&a.Function{Name: "=", Exec: eq})
	registerFunction(&a.Function{Name: "!=", Exec: neq})
	registerFunction(&a.Function{Name: ">", Exec: gt})
	registerFunction(&a.Function{Name: ">=", Exec: gte})
	registerFunction(&a.Function{Name: "<", Exec: lt})
	registerFunction(&a.Function{Name: "<=", Exec: lte})
}
