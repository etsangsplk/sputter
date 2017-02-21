package builtins

import a "github.com/kode4food/sputter/api"

func define(n a.Name, argNames a.Sequence, body a.Sequence) *a.Function {
	ac := a.Count(argNames)

	return &a.Function{
		Name: n,
		Exec: func(c a.Context, args a.Sequence) a.Value {
			a.AssertArity(args, ac)
			l := a.ChildContext(c)
			anIter := argNames.Iterate()
			aIter := args.Iterate()
			for ns, nok := anIter.Next(); nok; {
				an := ns.(*a.Symbol).Name
				av, _ := aIter.Next()
				l.Put(an, a.Eval(c, av))
				ns, nok = anIter.Next()
			}
			return a.EvalSequence(l, body)
		},
	}
}

func defun(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 3)
	g := c.Globals()
	i := args.Iterate()

	fv, _ := i.Next()
	fn := fv.(*a.Symbol).Name

	av, _ := i.Next()
	an := av.(a.Sequence)

	b := i.Rest()

	d := define(fn, an, b)
	a.PutFunction(g, d)
	return d
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "defun", Exec: defun})
}
