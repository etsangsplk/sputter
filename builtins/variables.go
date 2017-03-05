package builtins

import a "github.com/kode4food/sputter/api"

func def(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	g := a.GetContextNamespace(c)

	i := args.Iterate()
	s, _ := i.Next()
	n := a.AssertUnqualified(s).Name
	_, b := g.Get(n)
	if !b {
		v, _ := i.Next()
		g.Put(n, a.Eval(c, v))
	}
	return s
}

func let(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	l := a.ChildContext(c)
	i := args.Iterate()
	b, _ := i.Next()

	bi := b.(a.Sequence).Iterate()
	for s, ok := bi.Next(); ok; s, ok = bi.Next() {
		n := a.AssertUnqualified(s).Name
		if v, ok := bi.Next(); ok {
			l.Put(n, a.Eval(l, v))
		}
	}

	return a.EvalSequence(l, i.Rest())
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "def", Apply: def})
	a.PutFunction(Context, &a.Function{Name: "let", Apply: let})
}
