package builtins

import a "github.com/kode4food/sputter/api"

func ns(c a.Context, args a.Sequence) a.Value {
	i := args.Iterate()
	v, _ := i.Next()
	s := a.AssertUnqualifiedSymbol(v)
	c.Put(a.ContextDomain, s.Name)
	return a.GetContextNamespace(c)
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "ns", Apply: ns})
}
