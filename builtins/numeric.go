package builtins

import a "github.com/kode4food/sputter/api"

const (
	posInfName = "inf"
	negInfName = "-inf"
	incName    = "inc"
	decName    = "dec"
	addName    = "+"
	subName    = "-"
	mulName    = "*"
	divName    = "/"
	modName    = "mod"
	eqName     = "="
	neqName    = "!="
	gtName     = ">"
	gteName    = ">="
	ltName     = "<"
	lteName    = "<="

	isPosInfName = "is-pos-inf"
	isNegInfName = "is-neg-inf"
	isNaNName    = "is-nan"
)

type (
	reduceFunc  func(prev a.Number, next a.Number) a.Number
	compareFunc func(prev a.Number, next a.Number) bool

	incFunction      struct{ BaseBuiltIn }
	decFunction      struct{ BaseBuiltIn }
	addFunction      struct{ BaseBuiltIn }
	subFunction      struct{ BaseBuiltIn }
	mulFunction      struct{ BaseBuiltIn }
	divFunction      struct{ BaseBuiltIn }
	modFunction      struct{ BaseBuiltIn }
	eqFunction       struct{ BaseBuiltIn }
	neqFunction      struct{ BaseBuiltIn }
	gtFunction       struct{ BaseBuiltIn }
	gteFunction      struct{ BaseBuiltIn }
	ltFunction       struct{ BaseBuiltIn }
	lteFunction      struct{ BaseBuiltIn }
	isPosInfFunction struct{ BaseBuiltIn }
	isNegInfFunction struct{ BaseBuiltIn }
	isNaNFunction    struct{ BaseBuiltIn }
)

func reduceNum(s a.Sequence, v a.Number, fn reduceFunc) a.Value {
	res := v
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		fv := f.(a.Number)
		res = fn(res, fv)
	}
	return res
}

func fetchFirstNumber(args a.Vector) (a.Number, a.Sequence) {
	a.AssertMinimumArity(args, 1)
	return args[0].(a.Number), args[1:]
}

func compare(_ a.Context, s a.Vector, fn compareFunc) a.Bool {
	cur, rest := fetchFirstNumber(s)
	for f, r, ok := rest.Split(); ok; f, r, ok = r.Split() {
		v := f.(a.Number)
		if !fn(cur, v) {
			return a.False
		}
		cur = v
	}
	return a.True
}

func (*incFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 1)
	nv := args[0].(a.Number)
	return nv.Add(a.One)
}

func (*decFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 1)
	nv := args[0].(a.Number)
	return nv.Sub(a.One)
}

func (*addFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if !args.IsSequence() {
		return a.Zero
	}
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Add(n)
	})
}

func (*subFunction) Apply(_ a.Context, args a.Vector) a.Value {
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Sub(n)
	})
}

func (*mulFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if !args.IsSequence() {
		return a.One
	}
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Mul(n)
	})
}

func (*divFunction) Apply(_ a.Context, args a.Vector) a.Value {
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Div(n)
	})
}

func (*modFunction) Apply(_ a.Context, args a.Vector) a.Value {
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Mod(n)
	})
}

func (*eqFunction) Apply(c a.Context, args a.Vector) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.EqualTo
	})
}

func (*neqFunction) Apply(c a.Context, args a.Vector) a.Value {
	return !compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.EqualTo
	})
}

func (*gtFunction) Apply(c a.Context, args a.Vector) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.GreaterThan
	})
}

func (*gteFunction) Apply(c a.Context, args a.Vector) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.GreaterThan
	})
}

func (*ltFunction) Apply(c a.Context, args a.Vector) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.LessThan
	})
}

func (*lteFunction) Apply(c a.Context, args a.Vector) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.LessThan
	})
}

func (*isPosInfFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if n, ok := args[0].(a.Number); ok {
		if a.PosInfinity.Cmp(n) == a.EqualTo {
			return a.True
		}
	}
	return a.False
}

func (*isNegInfFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if n, ok := args[0].(a.Number); ok {
		if a.NegInfinity.Cmp(n) == a.EqualTo {
			return a.True
		}
	}
	return a.False
}

func (*isNaNFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if n, ok := args[0].(a.Number); ok {
		if n.IsNaN() {
			return a.True
		}
		return a.False
	}
	return a.True
}

func init() {
	var inc *incFunction
	var dec *decFunction
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
	var isPosInf *isPosInfFunction
	var isNegInf *isNegInfFunction
	var isNaN *isNaNFunction

	Namespace.Put(posInfName, a.PosInfinity)
	Namespace.Put(negInfName, a.NegInfinity)

	RegisterBuiltIn(incName, inc)
	RegisterBuiltIn(decName, dec)
	RegisterBuiltIn(addName, add)
	RegisterBuiltIn(subName, sub)
	RegisterBuiltIn(mulName, mul)
	RegisterBuiltIn(divName, div)
	RegisterBuiltIn(modName, mod)
	RegisterBuiltIn(eqName, eq)
	RegisterBuiltIn(neqName, neq)
	RegisterBuiltIn(gtName, gt)
	RegisterBuiltIn(gteName, gte)
	RegisterBuiltIn(ltName, lt)
	RegisterBuiltIn(lteName, lte)
	RegisterBuiltIn(isPosInfName, isPosInf)
	RegisterBuiltIn(isNegInfName, isNegInf)
	RegisterBuiltIn(isNaNName, isNaN)
}
