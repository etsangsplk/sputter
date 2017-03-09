package builtins

import a "github.com/kode4food/sputter/api"

func withNamespace(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)

	n := a.AssertUnqualified(args.First()).Name
	ns := a.GetNamespace(n)

	sc := a.ChildContext(c)
	sc.Put(a.ContextDomain, ns)
	return a.EvalSequence(sc, args.Rest())
}

func getNamespace(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name
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