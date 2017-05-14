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

var (
	emptyMetadata = a.Metadata{}
	defaultName   = a.Name("<lambda>")
)

func argNames(n a.Sequence) []a.Name {
	an := []a.Name{}
	for i := n; i.IsSequence(); i = i.Rest() {
		v := a.AssertUnqualified(i.First()).Name()
		an = append(an, v)
	}
	return an
}

func optionalMetadata(c a.Context, args a.Sequence) (a.Metadata, a.Sequence) {
	r := args
	var md a.Metadata
	if s, ok := r.First().(a.Str); ok {
		md = a.Metadata{a.MetaDoc: s}
		r = r.Rest()
	} else {
		md = emptyMetadata
	}

	if m, ok := r.First().(a.Mapped); ok {
		em := a.Eval(c, m).(a.Mapped)
		md = md.Merge(toMetadata(em))
		r = r.Rest()
	}
	return md, r
}

func optionalName(args a.Sequence) (a.Name, a.Sequence) {
	f := args.First()
	if s, ok := f.(a.Symbol); ok {
		if s.Domain() == a.LocalDomain {
			return s.Name(), args.Rest()
		}
		panic(a.Err(a.ExpectedUnqualified, s.Qualified()))
	}
	return defaultName, args
}

func getFunctionDefinition(c a.Context, args a.Sequence) *functionDefinition {
	a.AssertMinimumArity(args, 3)
	fn := a.AssertUnqualified(args.First()).Name()
	md, r := optionalMetadata(c, args.Rest())
	an := a.AssertVector(r.First())

	return &functionDefinition{
		args: an,
		body: r.Rest(),
		meta: md.Merge(a.Metadata{
			a.MetaName: fn,
			a.MetaArgs: an,
		}),
	}
}

func defineFunction(closure a.Context, d *functionDefinition) *a.Function {
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
	}).WithMetadata(d.meta).(*a.Function)
}

func defn(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(c, args)
	f := defineFunction(c, fd)
	a.GetContextNamespace(c).Put(f.Name(), f)
	return f
}

func fn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	fn, r := optionalName(args)
	md, r := optionalMetadata(c, r)
	an := a.AssertVector(r.First())

	return defineFunction(c, &functionDefinition{
		args: an,
		body: r.Rest(),
		meta: md.Merge(a.Metadata{
			a.MetaName: fn,
			a.MetaArgs: an,
		}),
	})
}

func apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f := a.AssertApplicable(a.Eval(c, args.First()))
	s := a.AssertSequence(a.Eval(c, args.Rest().First()))
	return f.Apply(c, s)
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
