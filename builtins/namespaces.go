package builtins

import a "github.com/kode4food/sputter/api"

func ns(c a.Context, args a.Sequence) a.Value {
	i := args.Iterate()
	v, _ := i.Next()
	n := a.AssertUnqualified(v).Name
	c.Put(a.ContextDomain, n)
	return a.GetContextNamespace(c)
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "ns", Apply: ns})
}
