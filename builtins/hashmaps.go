package builtins

import a "github.com/kode4food/sputter/api"

func hashMap(c a.Context, args a.Sequence) a.Value {
	if cnt, ok := args.(a.Countable); ok {
		l := cnt.Count()
		if l%2 != 0 {
			panic(a.ExpectedPair)
		}
		ml := l / 2
		r := make(a.ArrayMap, ml)
		i := a.Iterate(args)
		for idx := 0; idx < ml; idx++ {
			k, _ := i.Next()
			v, _ := i.Next()
			r[idx] = a.Vector{
				a.Eval(c, k),
				a.Eval(c, v),
			}
		}
		return r
	}

	r := a.ArrayMap{}
	i := a.Iterate(args)
	for k, ok := i.Next(); ok; k, ok = i.Next() {
		if v, ok := i.Next(); ok {
			r = append(r, a.Vector{
				a.Eval(c, k),
				a.Eval(c, v),
			})
		} else {
			panic(a.ExpectedPair)
		}
	}
	return r
}

func toHashMap(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	arg := a.Eval(c, args.First())
	seq := a.AssertSequence(arg)
	return hashMap(c, seq)
}

func isHashMap(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	if _, ok := a.Eval(c, v).(a.Mapped); ok {
		return a.True
	}
	return a.False
}

func init() {
	registerFunction(&a.Function{Name: "hash-map", Exec: hashMap})
	registerFunction(&a.Function{Name: "to-hash-map", Exec: toHashMap})
	registerPredicate(&a.Function{Name: "hash-map?", Exec: isHashMap})
}
