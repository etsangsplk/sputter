package builtins

import a "github.com/kode4food/sputter/api"

func identical(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	i := args.Iterate()
	l, _ := i.Next()
	r, _ := i.Next()
	if a.Eval(c, l) == a.Eval(c, r) {
		return a.True
	}
	return a.False
}
func isNil(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	i := args.Iterate()
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		if a.Eval(c, v) != a.Nil {
			return a.False
		}
	}
	return a.True
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "eq", Exec: identical})
	a.PutFunction(Context, &a.Function{Name: "nil?", Exec: isNil})
}
