package builtins

import a "github.com/kode4food/sputter/api"

type functionDefinition struct {
	name     a.Name
	doc      string
	argNames a.Sequence
	body     a.Sequence
}

func argNames(n a.Sequence) []a.Name {
	an := []a.Name{}
	for i := n; i.IsSequence(); i = i.Rest() {
		v := a.AssertUnqualified(i.First()).Name()
		an = append(an, v)
	}
	return an
}

func getFunctionDefinition(args a.Sequence) *functionDefinition {
	a.AssertMinimumArity(args, 3)

	i := a.Iterate(args)
	fv, _ := i.Next()
	fn := a.AssertUnqualified(fv).Name()

	var ds string
	av, _ := i.Next()
	if vs, ok := av.(string); ok {
		ds = vs
		av, _ = i.Next()
	}
	an := a.AssertSequence(av)

	return &functionDefinition{
		name:     fn,
		doc:      ds,
		argNames: an,
		body:     i.Rest(),
	}
}

func defineFunction(closure a.Context, d *functionDefinition) a.Function {
	an := argNames(d.argNames)
	ac := len(an)
	db := d.body

	return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		a.AssertArity(args, ac)
		l := a.ChildContext(closure)
		i := a.Iterate(args)
		for _, n := range an {
			v, _ := i.Next()
			l.Put(n, a.Eval(c, v))
		}
		return a.EvalSequence(l, db)
	}).WithMetadata(a.Metadata{
		a.MetaName: d.name,
		a.MetaDoc:  d.doc,
	}).(a.Function)
}

func defn(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(args)
	f := defineFunction(c, fd)
	a.GetContextNamespace(c).Put(fd.name, f)
	return f
}

func fn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	an := a.AssertSequence(args.First())

	return defineFunction(c, &functionDefinition{
		argNames: an,
		body:     args.Rest(),
	})
}

func apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f := a.AssertApplicable(a.Eval(c, args.First()))
	a := a.AssertSequence(a.Eval(c, args.Rest().First()))
	return f.Apply(c, a)
}

func init() {
	registerAnnotated(
		a.NewFunction(defn).WithMetadata(a.Metadata{
			a.MetaName: a.Name("defn"),
		}),
	)

	registerAnnotated(
		a.NewFunction(fn).WithMetadata(a.Metadata{
			a.MetaName: a.Name("fn"),
		}),
	)

	registerAnnotated(
		a.NewFunction(apply).WithMetadata(a.Metadata{
			a.MetaName: a.Name("apply"),
		}),
	)
}
