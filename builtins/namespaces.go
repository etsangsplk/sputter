package builtins

import a "github.com/kode4food/sputter/api"

const (
	withNamespaceName = "with-ns"
	getNamespaceName  = "ns"
)

type (
	withNamespaceFunction struct{ BaseBuiltIn }
	getNamespaceFunction  struct{ BaseBuiltIn }
)

func (*withNamespaceFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)

	f, r, _ := args.Split()
	n := a.AssertUnqualified(f).Name()
	ns := a.GetNamespace(n)

	sc := a.WithNamespace(a.ChildContext(c), ns)
	sc.Put(a.ContextDomain, ns)
	return a.MakeBlock(r).Eval(sc)
}

func (*getNamespaceFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name()
	return a.GetNamespace(n)
}

func init() {
	var withNamespace *withNamespaceFunction
	var getNamespace *getNamespaceFunction

	RegisterBuiltIn(withNamespaceName, withNamespace)
	RegisterBuiltIn(getNamespaceName, getNamespace)
}
