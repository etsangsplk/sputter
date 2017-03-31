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

func isVector(v a.Value) bool {
	if _, ok := v.(a.Vector); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(vector).WithMetadata(a.Metadata{
			a.MetaName: a.Name("vector"),
		}),
	)

	registerAnnotated(
		a.NewFunction(toVector).WithMetadata(a.Metadata{
			a.MetaName: a.Name("to-vector"),
		}),
	)

	registerSequencePredicate(isVector, a.Metadata{
		a.MetaName: a.Name("vector?"),
	})
}
