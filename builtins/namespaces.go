package builtins

import a "github.com/kode4food/sputter/api"

func withNamespace(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	i := args.Iterate()

	v, _ := i.Next()
	n := a.AssertUnqualified(v).Name
	ns := a.GetNamespace(n)

	sc := a.ChildContext(c)
	sc.Put(a.ContextDomain, ns)
	return a.EvalSequence(sc, i.Rest())
}

func getNamespace(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	i := args.Iterate()

	v, _ := i.Next()
	n := a.AssertUnqualified(v).Name
	return a.GetNamespace(n)
}

func init() {
	registerFunction(&a.Function{
		Name:  "with-ns",
		Apply: withNamespace,
	})

	registerFunction(&a.Function{
		Name:  "ns",
		Apply: getNamespace,
	})
}
