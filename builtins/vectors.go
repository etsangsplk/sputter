package builtins

import a "github.com/kode4food/sputter/api"

func vector(c a.Context, args a.Sequence) a.Value {
	l := a.Count(args)
	r := make(a.Vector, l)
	i := a.Iterate(args)
	for idx := 0; idx < l; idx++ {
		v, _ := i.Next()
		r[idx] = a.Eval(c, v)
	}
	return r
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
	registerFunction(&a.Function{Name: "vector", Apply: vector})
	registerPredicate(&a.Function{Name: "vector?", Apply: isVector})
}