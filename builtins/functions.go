package builtins

import a "github.com/kode4food/sputter/api"

type functionDefinition struct {
	name     a.Name
	doc      string
	argNames a.Sequence
	body     a.Sequence
	closure  a.Context
}

func argNames(n a.Sequence) []a.Name {
	an := []a.Name{}
	i := a.Iterate(n)
	for e, ok := i.Next(); ok; e, ok = i.Next() {
		v := a.AssertUnqualified(e).Name
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
		Doc:  d.doc,
		Exec: func(c a.Context, args a.Sequence) a.Value {
			a.AssertArity(args, ac)
			l := a.ChildContext(dc)
			i := a.Iterate(args)
			for _, n := range an {
				v, _ := i.Next()
				l.Put(n, a.Eval(c, v))
			}
			return a.EvalSequence(l, db)
		},
	}
}

func defn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 3)
	ns := a.GetContextNamespace(c)

	i := a.Iterate(args)
	fv, _ := i.Next()
	fn := a.AssertUnqualified(fv).Name

	var ds string
	av, _ := i.Next()
	if vs, ok := av.(string); ok {
		ds = vs
		av, _ = i.Next()
	}
	an := a.AssertSequence(av)

	d := define(&functionDefinition{
		name:     fn,
		doc:      ds,
		argNames: an,
		body:     i.Rest(),
		closure:  c,
	})

	ns.Put(fn, d)
	return d
}

func fn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	i := a.Iterate(args)
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
	f := a.AssertApplicable(a.Eval(c, args.First()))
	a := a.AssertSequence(a.Eval(c, args.Rest().First()))
	return f.Apply(c, a)
}

func init() {
	registerFunction(&a.Function{Name: "defn", Exec: defn})
	registerFunction(&a.Function{Name: "fn", Exec: fn})
	registerFunction(&a.Function{Name: "apply", Exec: apply})
}
