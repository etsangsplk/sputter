package builtins

import a "github.com/kode4food/sputter/api"

func define(n a.Name, argNames a.Sequence, body a.Sequence) *a.Function {
	ac := argNames.Count()

	return &a.Function{
		Name: n,
		Exec: func(c *a.Context, args a.Sequence) a.Value {
			AssertArity(args, ac)
			l := c.Child()
			anIter := argNames.Iterate()
			aIter := args.Iterate()
			for ns, nok := anIter.Next(); nok; {
				an := ns.(*a.Symbol).Name
				av, aok := aIter.Next()
				if aok {
					l.Put(an, av)
				} else {
					l.Put(an, a.Nil)
				}
				ns, nok = anIter.Next()
			}
			return a.EvalIterator(l, body.Iterate())
		},
	}
}

func defun(c *a.Context, args a.Sequence) a.Value {
	AssertMinimumArity(args, 3)
	g := c.Globals()
	i := args.Iterate()

	fv, _ := i.Next()
	fn := fv.(*a.Symbol).Name

	av, _ := i.Next()
	an := av.(a.Sequence)

	b := i.Iterable()

	d := define(fn, an, b)
	g.PutFunction(d)
	return d
}

func init() {
	Context.PutFunction(&a.Function{Name: "defun", Exec: defun})
}
