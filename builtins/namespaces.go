package builtins

import a "github.com/kode4food/sputter/api"

func ns(c a.Context, args a.Sequence) a.Value {
	i := args.Iterate()
	v, _ := i.Next()
	n := a.AssertUnqualified(v).Name
	ns := a.GetNamespace(n)
	c.Put(a.ContextDomain, ns)
	return ns
}

func init() {
	putFunction(BuiltInNamespace, &a.Function{Name: "ns", Apply: ns})
}
