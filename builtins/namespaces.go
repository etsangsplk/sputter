package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func withNamespace(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)

	n := a.AssertUnqualified(args.First()).Name()
	ns := a.GetNamespace(n)

	sc := a.WithNamespace(a.ChildContext(c), ns)
	sc.Put(a.ContextDomain, ns)
	return a.EvalSequence(sc, args.Rest())
}

func getNamespace(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name()
	return a.GetNamespace(n)
}

func init() {
	registerAnnotated(
		a.NewFunction(withNamespace).WithMetadata(a.Metadata{
			a.MetaName: a.Name("with-ns"),
			a.MetaDoc:  d.Get("with-ns"),
		}),
	)

	registerAnnotated(
		a.NewFunction(getNamespace).WithMetadata(a.Metadata{
			a.MetaName: a.Name("ns"),
			a.MetaDoc:  d.Get("ns"),
		}),
	)
}
