package builtins

import a "github.com/kode4food/sputter/api"

type (
	reduceFunc  func(prev a.Number, next a.Number) a.Number
	compareFunc func(prev a.Number, next a.Number) bool
)

func reduceNum(s a.Sequence, v a.Number, f reduceFunc) a.Value {
	r := v
	for i := s; i.IsSequence(); i = i.Rest() {
		fv := a.AssertNumber(i.First())
		r = f(r, fv)
	}
	return r
}

func fetchFirstNumber(args a.Sequence) (a.Number, a.Sequence) {
	a.AssertMinimumArity(args, 1)
	nv := a.AssertNumber(args.First())
	return nv, args.Rest()
}

func compare(_ a.Context, s a.Sequence, f compareFunc) a.Value {
	cur, r := fetchFirstNumber(s)
	for i := r; i.IsSequence(); i = i.Rest() {
		v := a.AssertNumber(i.First())
		if !f(cur, v) {
			return a.False
		}
		cur = v
	}
	return a.True
}

func add(_ a.Context, args a.Sequence) a.Value {
	if !args.IsSequence() {
		return a.Zero
	}
	f, r := fetchFirstNumber(args)
	return reduceNum(r, f, func(p a.Number, n a.Number) a.Number {
		return p.Add(n)
	})
}

func sub(_ a.Context, args a.Sequence) a.Value {
	f, r := fetchFirstNumber(args)
	return reduceNum(r, f, func(p a.Number, n a.Number) a.Number {
		return p.Sub(n)
	})
}

func mul(_ a.Context, args a.Sequence) a.Value {
	if !args.IsSequence() {
		return a.One
	}
	f, r := fetchFirstNumber(args)
	return reduceNum(r, f, func(p a.Number, n a.Number) a.Number {
		return p.Mul(n)
	})
}

func div(_ a.Context, args a.Sequence) a.Value {
	f, r := fetchFirstNumber(args)
	return reduceNum(r, f, func(p a.Number, n a.Number) a.Number {
		return p.Div(n)
	})
}

func mod(_ a.Context, args a.Sequence) a.Value {
	f, r := fetchFirstNumber(args)
	return reduceNum(r, f, func(p a.Number, n a.Number) a.Number {
		return p.Mod(n)
	})
}

func eq(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
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
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.GreaterThan
	})
}

func gte(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.GreaterThan
	})
}

func lt(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.LessThan
	})
}

func lte(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.LessThan
	})
}

func init() {
	Namespace.Put("inf", a.PosInfinity)
	Namespace.Put("-inf", a.NegInfinity)

	RegisterBuiltIn("+", add)
	RegisterBuiltIn("-", sub)
	RegisterBuiltIn("*", mul)
	RegisterBuiltIn("/", div)
	RegisterBuiltIn("%", mod)
	RegisterBuiltIn("=", eq)
	RegisterBuiltIn("!=", neq)
	RegisterBuiltIn(">", gt)
	RegisterBuiltIn(">=", gte)
	RegisterBuiltIn("<", lt)
	RegisterBuiltIn("<=", lte)
}
