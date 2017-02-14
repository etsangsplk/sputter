package builtins

import a "github.com/kode4food/sputter/api"

func defvar(c *a.Context, args a.Iterable) a.Value {
	AssertMinimumArity(args, 2)
	g := c.Globals()
	i := args.Iterate()
	s, _ := i.Next()
	n := s.(*a.Symbol).Name
	_, b := g.Get(n)
	if !b {
		v, _ := i.Next()
		g.Put(n, a.Evaluate(c, v))
	}
	return s
}

func let(c *a.Context, args a.Iterable) a.Value {
	AssertMinimumArity(args, 2)
	l := c.Child()
	i := args.Iterate()
	b, _ := i.Next()

	bi := b.(a.Iterable).Iterate()
	for s, ok := bi.Next(); ok; s, ok = bi.Next() {
		n := s.(*a.Symbol).Name
		if v, ok := bi.Next(); ok {
			l.Put(n, a.Evaluate(l, v))
		}
	}

	return a.EvaluateIterator(l, i)
}

func init() {
	Context.PutFunction(&a.Function{Name: "defvar", Exec: defvar})
	Context.PutFunction(&a.Function{Name: "let", Exec: let})
}
