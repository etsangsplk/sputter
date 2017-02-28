package builtins

import a "github.com/kode4food/sputter/api"

func defvar(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	g := a.GetNamespace(a.UserDomain)

	i := args.Iterate()
	s, _ := i.Next()
	n := s.(*a.Symbol).Name
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
		n := s.(*a.Symbol).Name
		if v, ok := bi.Next(); ok {
			l.Put(n, a.Eval(l, v))
		}
	}

	return a.EvalSequence(l, i.Rest())
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "defvar", Exec: defvar})
	a.PutFunction(Context, &a.Function{Name: "let", Exec: let})
}
