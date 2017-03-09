package builtins

import a "github.com/kode4food/sputter/api"

// BuiltIns is a special Namespace for built-in identifiers
var BuiltIns = a.GetNamespace(a.BuiltInDomain)

func registerFunction(f *a.Function) {
	BuiltIns.Put(f.Name, f)
}

func registerPredicate(f *a.Function) {
	registerFunction(f)
	registerFunction(&a.Function{
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

func do(c a.Context, args a.Sequence) a.Value {
	return a.EvalSequence(c, args)
}

func quote(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	return args.First()
}

func init() {
	registerFunction(&a.Function{Name: "do", Apply: do})

	registerFunction(&a.Function{
		Name:  "quote",
		Apply: quote,
		Data:  true,
	})
}
