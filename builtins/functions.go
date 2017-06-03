package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

// InvalidRestArgument is thrown if you include more than one rest argument
const InvalidRestArgument = "rest-argument not well-formed: %s"

type functionDefinition struct {
	args a.Vector
	rest bool
	body a.Sequence
	meta a.Metadata
}

type argProcessor func(c a.Context, args a.Sequence) a.Context

var (
	emptyMetadata = a.Metadata{}
	defaultName   = a.Name("<lambda>")
	restMarker    = a.Name("&")
	metaDocAsset  = a.NewKeyword("doc-asset")
)

func makeArgProcessor(cl a.Context, s a.Sequence) argProcessor {
	an := []a.Name{}
	for i := s; i.IsSequence(); i = i.Rest() {
		n := a.AssertUnqualified(i.First()).Name()
		if n == restMarker {
			rn := parseRestArg(i)
			return makeRestArgProcessor(cl, an, rn)
		}
		an = append(an, n)
	}
	return makeFixedArgProcessor(cl, an)
}

func parseRestArg(s a.Sequence) a.Name {
	r := s.Rest()
	if r.IsSequence() {
		n := a.AssertUnqualified(r.First()).Name()
		if n != restMarker && !r.Rest().IsSequence() {
			return n
		}
	}
	panic(a.Err(InvalidRestArgument, s))
}

func makeRestArgProcessor(cl a.Context, an []a.Name, rn a.Name) argProcessor {
	ac := len(an)

	return func(c a.Context, args a.Sequence) a.Context {
		a.AssertMinimumArity(args, ac)
		l := a.ChildContext(cl)
		i := args
		for _, n := range an {
			l.Put(n, i.First().Eval(c))
			i = i.Rest()
		}

		r := []a.Value{}
		for ; i.IsSequence(); i = i.Rest() {
			r = append(r, i.First().Eval(c))
		}
		l.Put(rn, a.NewList(r...))
		return l
	}
}

func makeFixedArgProcessor(cl a.Context, an []a.Name) argProcessor {
	ac := len(an)

	return func(c a.Context, args a.Sequence) a.Context {
		a.AssertArity(args, ac)
		l := a.ChildContext(cl)
		i := args
		for _, n := range an {
			l.Put(n, i.First().Eval(c))
			i = i.Rest()
		}
		return l
	}
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
		em := m.Eval(c).(a.Mapped)
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

func loadDocumentation(md a.Metadata) a.Metadata {
	v, ok := md.Get(metaDocAsset)
	if !ok {
		return md
	}

	fn, ok := v.(a.Str)
	if !ok || !fn.IsSequence() {
		return md
	}

	s := string(fn)
	if !d.Exists(s) {
		return md
	}

	return md.Merge(a.Metadata{
		a.MetaDoc: d.Get(s),
	})
}

func getFunctionDefinition(c a.Context, args a.Sequence) *functionDefinition {
	a.AssertMinimumArity(args, 3)
	fn := a.AssertUnqualified(args.First()).Name()
	md, r := optionalMetadata(c, args.Rest())
	md = loadDocumentation(md)
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

func defineFunction(closure a.Context, d *functionDefinition) a.Function {
	ap := makeArgProcessor(closure, d.args)
	db := d.body

	return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		l := ap(c, args)
		return a.EvalBlock(l, db)
	}).WithMetadata(d.meta).(a.Function)
}

func makefn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	fn, r := optionalName(args)
	md, r := optionalMetadata(c, r)
	md = loadDocumentation(md)
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
	f := a.AssertApplicable(args.First().Eval(c))
	s := a.AssertSequence(args.Rest().First().Eval(c))
	return f.Apply(c, s)
}

func init() {
	registerAnnotated(
		a.NewFunction(makefn).WithMetadata(a.Metadata{
			a.MetaName: a.Name("make-fn"),
		}),
	)

	registerAnnotated(
		a.NewFunction(apply).WithMetadata(a.Metadata{
			a.MetaName: a.Name("apply"),
			a.MetaDoc:  d.Get("apply"),
		}),
	)
}
