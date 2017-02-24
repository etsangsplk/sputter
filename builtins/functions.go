package builtins

import a "github.com/kode4food/sputter/api"

type functionDefinition struct {
	name     a.Name
	argNames a.Sequence
	body     a.Sequence
	closure  a.Context
}

func define(d *functionDefinition) *a.Function {
	ac := a.Count(d.argNames)

	return &a.Function{
		Name: d.name,
		Exec: func(c a.Context, args a.Sequence) a.Value {
			a.AssertArity(args, ac)
			l := a.ChildContext(d.closure)
			anIter := d.argNames.Iterate()
			aIter := args.Iterate()
			for ns, nok := anIter.Next(); nok; {
				an := ns.(*a.Symbol).Name
				av, _ := aIter.Next()
				l.Put(an, a.Eval(c, av))
				ns, nok = anIter.Next()
			}
			return a.EvalSequence(l, d.body)
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

	d := define(&functionDefinition{
		name:     fn,
		argNames: an,
		body:     b,
		closure:  c,
	})

	a.PutFunction(g, d)
	return d
}

func lambda(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	i := args.Iterate()
	av, _ := i.Next()
	an := av.(a.Sequence)

	b := i.Rest()

	return define(&functionDefinition{
		argNames: an,
		body:     b,
		closure:  c,
	})
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "defun", Exec: defun})
	a.PutFunction(Context, &a.Function{Name: "lambda", Exec: lambda})
}
