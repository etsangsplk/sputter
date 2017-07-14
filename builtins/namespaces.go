package builtins

import a "github.com/kode4food/sputter/api"

func withNamespace(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)

	n := a.AssertUnqualified(args.First()).Name()
	ns := a.GetNamespace(n)

	sc := a.WithNamespace(a.ChildContext(c), ns)
	sc.Put(a.ContextDomain, ns)
	return a.EvalBlock(sc, args.Rest())
}

func getNamespace(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name()
	return a.GetNamespace(n)
}

func init() {
	RegisterBuiltIn("with-ns", withNamespace)
	RegisterBuiltIn("ns", getNamespace)
}
