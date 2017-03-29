package builtins

import a "github.com/kode4food/sputter/api"

func def(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	ns := a.GetContextNamespace(c)

	s := args.First()
	n := a.AssertUnqualified(s).Name
	v := args.Rest().First()
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
	registerAnnotated(
		a.NewFunction(def).WithMetadata(a.Metadata{
			a.MetaName: a.Name("def"),
		}),
	)

	registerAnnotated(
		a.NewFunction(let).WithMetadata(a.Metadata{
			a.MetaName: a.Name("let"),
		}),
	)
}
