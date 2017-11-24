package builtins

import a "github.com/kode4food/sputter/api"

const (
	withNamespaceName = "with-ns"
	getNamespaceName  = "ns"
	namespacePutName  = "ns-put"
)

type (
	withNamespaceFunction struct{ BaseBuiltIn }
	getNamespaceFunction  struct{ BaseBuiltIn }
	namespacePutFunction  struct{ BaseBuiltIn }
)

func (*withNamespaceFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)

	f, r, _ := args.Split()
	n := f.(a.LocalSymbol).Name()
	ns := a.GetNamespace(n)

	sc := a.WithNamespace(a.ChildContext(c), ns)
	sc.Put(a.ContextDomain, ns)
	return a.MakeBlock(r).Eval(sc)
}

func (*getNamespaceFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := args.First().(a.LocalSymbol).Name()
	return a.GetNamespace(n)
}

func (*namespacePutFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 3)

	f, r, _ := args.Split()
	ns := f.(a.Namespace)

	f, r, _ = r.Split()
	n := f.(a.LocalSymbol).Name()
	ns.Put(n, r.First())
	return f
}

func init() {
	var withNamespace *withNamespaceFunction
	var getNamespace *getNamespaceFunction
	var namespacePut *namespacePutFunction

	RegisterBuiltIn(withNamespaceName, withNamespace)
	RegisterBuiltIn(getNamespaceName, getNamespace)
	RegisterBuiltIn(namespacePutName, namespacePut)
}
