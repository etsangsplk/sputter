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

func (*withNamespaceFunction) Apply(c a.Context, args a.Values) a.Value {
	a.AssertMinimumArity(args, 2)

	n := args[0].(a.LocalSymbol).Name()
	ns := a.GetNamespace(n)

	lc := a.ChildLocals(c)
	sc := a.WithNamespace(lc, ns)
	sc.Put(a.ContextDomain, ns)
	return a.MakeBlock(args[1:]).Eval(sc)
}

func (*getNamespaceFunction) Apply(c a.Context, args a.Values) a.Value {
	if a.AssertArityRange(args, 0, 1) == 0 {
		return a.GetContextNamespace(c)
	}
	n := args[0].(a.LocalSymbol).Name()
	return a.GetNamespace(n)
}

func (*namespacePutFunction) Apply(c a.Context, args a.Values) a.Value {
	a.AssertArity(args, 3)

	ns := a.Eval(c, args[0]).(a.Namespace)

	n := args[1].(a.LocalSymbol).Name()
	ns.Put(n, a.Eval(c, args[2]))
	return n
}

func init() {
	var withNamespace *withNamespaceFunction
	var getNamespace *getNamespaceFunction
	var namespacePut *namespacePutFunction

	RegisterBuiltIn(withNamespaceName, withNamespace)
	RegisterBuiltIn(getNamespaceName, getNamespace)
	RegisterBuiltIn(namespacePutName, namespacePut)
}
