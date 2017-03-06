package builtins

import a "github.com/kode4food/sputter/api"

func registerPredicate(f *a.Function) {
	putFunction(BuiltInNamespace, f)
	putFunction(BuiltInNamespace, &a.Function{
		Name: "!" + f.Name,
		Apply: func(c a.Context, args a.Sequence) a.Value {
			r := f.Apply(c, args)
			if r == a.True {
				return a.False
			}
			return a.True
		},
	})
}

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
	registerPredicate(&a.Function{Name: "eq", Apply: identical})
	registerPredicate(&a.Function{Name: "nil?", Apply: isNil})
}
