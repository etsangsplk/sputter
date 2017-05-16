package builtins

import a "github.com/kode4food/sputter/api"

type reduceFunc func(prev *a.Number, next *a.Number) *a.Number
type compareFunc func(prev *a.Number, next *a.Number) bool

var (
	zero = a.NewFloat(0)
	one  = a.NewFloat(1)
)

func reduceNum(c a.Context, s a.Sequence, v *a.Number, f reduceFunc) a.Value {
	r := v
	for i := s; i.IsSequence(); i = i.Rest() {
		fv := a.AssertNumber(i.First().Eval(c))
		r = f(r, fv)
	}
	return r
}

func fetchFirstNumber(c a.Context, args a.Sequence) (*a.Number, a.Sequence) {
	a.AssertMinimumArity(args, 1)
	nv := a.AssertNumber(args.First().Eval(c))
	return nv, args.Rest()
}

func compare(c a.Context, s a.Sequence, f compareFunc) a.Value {
	cur, r := fetchFirstNumber(c, s)
	for i := r; i.IsSequence(); i = i.Rest() {
		v := a.AssertNumber(i.First().Eval(c))
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
	return reduceNum(c, r, f, func(p *a.Number, n *a.Number) *a.Number {
		return p.Add(n)
	})
}

func sub(c a.Context, args a.Sequence) a.Value {
	f, r := fetchFirstNumber(c, args)
	return reduceNum(c, r, f, func(p *a.Number, n *a.Number) *a.Number {
		return p.Sub(n)
	})
}

func mul(c a.Context, args a.Sequence) a.Value {
	if !args.IsSequence() {
		return one
	}
	f, r := fetchFirstNumber(c, args)
	return reduceNum(c, r, f, func(p *a.Number, n *a.Number) *a.Number {
		return p.Mul(n)
	})
}

func div(c a.Context, args a.Sequence) a.Value {
	f, r := fetchFirstNumber(c, args)
	return reduceNum(c, r, f, func(p *a.Number, n *a.Number) *a.Number {
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
	registerAnnotated(
		a.NewFunction(add).WithMetadata(a.Metadata{
			a.MetaName: a.Name("+"),
		}),
	)

	registerAnnotated(
		a.NewFunction(sub).WithMetadata(a.Metadata{
			a.MetaName: a.Name("-"),
		}),
	)

	registerAnnotated(
		a.NewFunction(mul).WithMetadata(a.Metadata{
			a.MetaName: a.Name("*"),
		}),
	)

	registerAnnotated(
		a.NewFunction(div).WithMetadata(a.Metadata{
			a.MetaName: a.Name("/"),
		}),
	)

	registerAnnotated(
		a.NewFunction(eq).WithMetadata(a.Metadata{
			a.MetaName: a.Name("="),
		}),
	)

	registerAnnotated(
		a.NewFunction(neq).WithMetadata(a.Metadata{
			a.MetaName: a.Name("!="),
		}),
	)

	registerAnnotated(
		a.NewFunction(gt).WithMetadata(a.Metadata{
			a.MetaName: a.Name(">"),
		}),
	)

	registerAnnotated(
		a.NewFunction(gte).WithMetadata(a.Metadata{
			a.MetaName: a.Name(">="),
		}),
	)

	registerAnnotated(
		a.NewFunction(lt).WithMetadata(a.Metadata{
			a.MetaName: a.Name("<"),
		}),
	)

	registerAnnotated(
		a.NewFunction(lte).WithMetadata(a.Metadata{
			a.MetaName: a.Name("<="),
		}),
	)
}
