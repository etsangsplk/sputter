package builtins

import a "github.com/kode4food/sputter/api"

func def(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	ns := a.GetContextNamespace(c)

	i := a.Iterate(args)
	s, _ := i.Next()
	n := a.AssertUnqualified(s).Name
	v, _ := i.Next()
	ns.Put(n, a.Eval(c, v))
	return s
}

func let(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	l := a.ChildContext(c)
	i := a.Iterate(args)
	b, _ := i.Next()

	bi := a.Iterate(a.AssertSequence(b))
	for s, ok := bi.Next(); ok; s, ok = bi.Next() {
		n := a.AssertUnqualified(s).Name
		if v, ok := bi.Next(); ok {
			l.Put(n, a.Eval(l, v))
		}
	}

	return a.EvalSequence(l, i.Rest())
}

func init() {
	registerFunction(&a.Function{Name: "def", Exec: def})
	registerFunction(&a.Function{Name: "let", Exec: let})
}
