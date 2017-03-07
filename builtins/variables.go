package builtins

import a "github.com/kode4food/sputter/api"

func def(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	ns := a.GetContextNamespace(c)

	i := args.Iterate()
	s, _ := i.Next()
	n := a.AssertUnqualified(s).Name
	v, _ := i.Next()
	ns.Put(n, a.Eval(c, v))
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
	putFunction(BuiltIns, &a.Function{Name: "def", Apply: def})
	putFunction(BuiltIns, &a.Function{Name: "let", Apply: let})
}
