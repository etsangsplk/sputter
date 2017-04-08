package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

type functionDefinition struct {
	args a.Vector
	body a.Sequence
	meta a.Metadata
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

	fn := a.AssertUnqualified(args.First()).Name()
	r := args.Rest()

	var ds string
	if vs, ok := r.First().(string); ok {
		ds = vs
		r = r.Rest()
	}
	an := a.AssertVector(r.First())

	return &functionDefinition{
		args: an,
		body: r.Rest(),
		meta: a.Metadata{
			a.MetaName: fn,
			a.MetaDoc:  ds,
		},
	}
}

func defineFunction(closure a.Context, d *functionDefinition) a.Function {
	an := argNames(d.args)
	ac := len(an)
	db := d.body

	return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		a.AssertArity(args, ac)
		l := a.ChildContext(closure)
		i := args
		for _, n := range an {
			l.Put(n, a.Eval(c, i.First()))
			i = i.Rest()
		}
		return a.EvalSequence(l, db)
	}).WithMetadata(d.meta).(a.Function)
}

func defn(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(args)
	f := defineFunction(c, fd)
	a.GetContextNamespace(c).Put(f.Name(), f)
	return f
}

func fn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	an := a.AssertVector(args.First())

	return defineFunction(c, &functionDefinition{
		args: an,
		body: args.Rest(),
		meta: a.Metadata{},
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
			a.MetaDoc:  d.Get("defn"),
		}),
	)

	registerAnnotated(
		a.NewFunction(fn).WithMetadata(a.Metadata{
			a.MetaName: a.Name("fn"),
			a.MetaDoc:  d.Get("fn"),
		}),
	)

	registerAnnotated(
		a.NewFunction(apply).WithMetadata(a.Metadata{
			a.MetaName: a.Name("apply"),
			a.MetaDoc:  d.Get("apply"),
		}),
	)
}
