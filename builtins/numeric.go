package builtins

import a "github.com/kode4food/sputter/api"

type (
	reduceFunc  func(prev a.Number, next a.Number) a.Number
	compareFunc func(prev a.Number, next a.Number) bool

	addFunction struct{ BaseBuiltIn }
	subFunction struct{ BaseBuiltIn }
	mulFunction struct{ BaseBuiltIn }
	divFunction struct{ BaseBuiltIn }
	modFunction struct{ BaseBuiltIn }
	eqFunction  struct{ BaseBuiltIn }
	neqFunction struct{ BaseBuiltIn }
	gtFunction  struct{ BaseBuiltIn }
	gteFunction struct{ BaseBuiltIn }
	ltFunction  struct{ BaseBuiltIn }
	lteFunction struct{ BaseBuiltIn }
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

func compare(_ a.Context, s a.Sequence, f compareFunc) a.Bool {
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

func (f *addFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if !args.IsSequence() {
		return a.Zero
	}
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Add(n)
	})
}

func (f *subFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Sub(n)
	})
}

func (f *mulFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if !args.IsSequence() {
		return a.One
	}
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Mul(n)
	})
}

func (f *divFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Div(n)
	})
}

func (f *modFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Mod(n)
	})
}

func (f *eqFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.EqualTo
	})
}

func (f *neqFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.EqualTo
	}).Not()
}

func (f *gtFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.GreaterThan
	})
}

func (f *gteFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.GreaterThan
	})
}

func (f *ltFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.LessThan
	})
}

func (f *lteFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.LessThan
	})
}

func isPosInfinity(v a.Value) bool {
	if n, ok := v.(a.Number); ok {
		return a.PosInfinity.Cmp(n) == a.EqualTo
	}
	return false
}

func isNegInfinity(v a.Value) bool {
	if n, ok := v.(a.Number); ok {
		return a.NegInfinity.Cmp(n) == a.EqualTo
	}
	return false
}

func isNaN(v a.Value) bool {
	if n, ok := v.(a.Number); ok {
		return n.IsNaN()
	}
	return true
}

func init() {
	var add *addFunction
	var sub *subFunction
	var mul *mulFunction
	var div *divFunction
	var mod *modFunction
	var eq *eqFunction
	var neq *neqFunction
	var gt *gtFunction
	var gte *gteFunction
	var lt *ltFunction
	var lte *lteFunction

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

	RegisterSequencePredicate("inf?", isPosInfinity)
	RegisterSequencePredicate("-inf?", isNegInfinity)
	RegisterSequencePredicate("nan?", isNaN)
}
