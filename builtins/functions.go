package builtins

import a "github.com/kode4food/sputter/api"

type functionDefinition struct {
	name     a.Name
	argNames a.Sequence
	body     a.Sequence
	closure  a.Context
}

func argNames(n a.Sequence) []a.Name {
	an := []a.Name{}
	i := n.Iterate()
	for e, ok := i.Next(); ok; e, ok = i.Next() {
		v := a.AssertSymbol(e).Name
		an = append(an, v)
	}
	return an
}

func define(d *functionDefinition) *a.Function {
	an := argNames(d.argNames)
	ac := len(an)
	dc := d.closure
	db := d.body

	return &a.Function{
		Name: d.name,
		Apply: func(c a.Context, args a.Sequence) a.Value {
			a.AssertArity(args, ac)
			l := a.ChildContext(dc)
			i := args.Iterate()
			for _, n := range an {
				v, _ := i.Next()
				l.Put(n, a.Eval(c, v))
			}
			return a.EvalSequence(l, db)
		},
	}
}

func defun(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 3)
	g := a.GetNamespace(a.UserDomain)

	i := args.Iterate()
	fv, _ := i.Next()
	fn := a.AssertSymbol(fv).Name
	av, _ := i.Next()
	an := a.AssertSequence(av)

	d := define(&functionDefinition{
		name:     fn,
		argNames: an,
		body:     i.Rest(),
		closure:  c,
	})

	a.PutFunction(g, d)
	return d
}

func lambda(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	i := args.Iterate()
	av, _ := i.Next()
	an := a.AssertSequence(av)

	return define(&functionDefinition{
		argNames: an,
		body:     i.Rest(),
		closure:  c,
	})
}

func apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	i := args.Iterate()
	fv, _ := i.Next()
	av, _ := i.Next()
	f := a.AssertFunction(a.Eval(c, fv))
	return f.Apply(c, a.AssertSequence(av))
	panic(a.ExpectedFunction)
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "defun", Apply: defun})
	a.PutFunction(Context, &a.Function{Name: "lambda", Apply: lambda})
	a.PutFunction(Context, &a.Function{Name: "apply", Apply: apply})
}
