package builtins

import a "github.com/kode4food/sputter/api"

func vector(c a.Context, args a.Sequence) a.Value {
	if cnt, ok := args.(a.Countable); ok {
		l := cnt.Count()
		r := make(a.Vector, l)
		idx := 0
		for i := args; i.IsSequence(); i = i.Rest() {
			r[idx] = a.Eval(c, i.First())
			idx++
		}
		return r
	}

	r := a.Vector{}
	for i := args; i.IsSequence(); i = i.Rest() {
		r = append(r, a.Eval(c, i.First()))
	}
	return r
}

func toVector(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	arg := a.Eval(c, args.First())
	seq := a.AssertSequence(arg)
	return vector(c, seq)
}

func isVector(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	if _, ok := a.Eval(c, v).(a.Vector); ok {
		return a.True
	}
	return a.False
}

func init() {
	registerAnnotated(
		a.NewFunction(vector).WithMetadata(a.Variables{
			a.MetaName: a.Name("vector"),
		}),
	)

	registerAnnotated(
		a.NewFunction(toVector).WithMetadata(a.Variables{
			a.MetaName: a.Name("to-vector"),
		}),
	)

	registerPredicate(
		a.NewFunction(isVector).WithMetadata(a.Variables{
			a.MetaName: a.Name("vector?"),
		}).(a.Function),
	)
}
