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

func init() {
	registerFunction(&a.Function{Name: "vector", Apply: vector})
}
