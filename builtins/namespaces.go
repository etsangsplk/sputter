package builtins

import a "github.com/kode4food/sputter/api"

type (
	withNamespaceFunction struct{ BaseBuiltIn }
	getNamespaceFunction  struct{ BaseBuiltIn }
)

func (f *withNamespaceFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)

	n := a.AssertUnqualified(args.First()).Name()
	ns := a.GetNamespace(n)

	sc := a.WithNamespace(a.ChildContext(c), ns)
	sc.Put(a.ContextDomain, ns)
	return a.MakeBlock(args.Rest()).Eval(sc)
}

func (f *getNamespaceFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name()
	return a.GetNamespace(n)
}

func init() {
	var withNamespace *withNamespaceFunction
	var getNamespace *getNamespaceFunction

	RegisterBuiltIn("with-ns", withNamespace)
	RegisterBuiltIn("ns", getNamespace)
}
