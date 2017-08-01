package builtins

import a "github.com/kode4food/sputter/api"

type (
	withNamespaceFunction struct{ a.ReflectedFunction }
	getNamespaceFunction  struct{ a.ReflectedFunction }
)

func (f *withNamespaceFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)

	n := a.AssertUnqualified(args.First()).Name()
	ns := a.GetNamespace(n)

	sc := a.WithNamespace(a.ChildContext(c), ns)
	sc.Put(a.ContextDomain, ns)
	return a.EvalBlock(sc, args.Rest())
}

func (f *getNamespaceFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name()
	return a.GetNamespace(n)
}

func init() {
	var withNamespace *withNamespaceFunction
	var getNamespace *getNamespaceFunction

	RegisterBaseFunction("with-ns", withNamespace)
	RegisterBaseFunction("ns", getNamespace)
}
