package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

// InvalidRestArgument is thrown if you include more than one rest argument
const InvalidRestArgument = "rest-argument not well-formed: %s"

type (
	functionDefinition struct {
		args a.Vector
		rest bool
		body a.Sequence
		meta a.Object
	}

	argProcessor func(a.Context, a.Sequence) a.Context
)

var (
	emptyMetadata = a.Properties{}
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
	panic(a.ErrStr(InvalidRestArgument, s))
}

func makeRestArgProcessor(cl a.Context, an []a.Name, rn a.Name) argProcessor {
	ac := len(an)

	return func(_ a.Context, args a.Sequence) a.Context {
		a.AssertMinimumArity(args, ac)
		l := a.ChildContext(cl)
		i := args
		for _, n := range an {
			l.Put(n, i.First())
			i = i.Rest()
		}
		l.Put(rn, a.ToList(i))
		return l
	}
}

func makeFixedArgProcessor(cl a.Context, an []a.Name) argProcessor {
	ac := len(an)

	return func(_ a.Context, args a.Sequence) a.Context {
		a.AssertArity(args, ac)
		l := a.ChildContext(cl)
		i := args
		for _, n := range an {
			l.Put(n, i.First())
			i = i.Rest()
		}
		return l
	}
}

func optionalMetadata(c a.Context, args a.Sequence) (a.Object, a.Sequence) {
	r := args
	var md a.Object
	if s, ok := r.First().(a.Str); ok {
		md = a.Properties{a.DocKey: s}
		r = r.Rest()
	} else {
		md = emptyMetadata
	}

	if m, ok := r.First().(a.MappedSequence); ok {
		em := a.Eval(c, m).(a.MappedSequence)
		md = md.Child(toProperties(em))
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
		panic(a.ErrStr(a.ExpectedUnqualified, s.Qualified()))
	}
	return defaultName, args
}

func loadDocumentation(md a.Object) a.Object {
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

	return md.Child(a.Properties{
		a.DocKey: d.Get(s),
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
		meta: md.Child(a.Properties{
			a.NameKey: fn,
			a.ArgsKey: an,
		}),
	}
}

func defineFunction(closure a.Context, d *functionDefinition) a.Function {
	ap := makeArgProcessor(closure, d.args)
	ex := a.MacroExpandAll(closure, d.body).(a.Sequence)
	db := a.NewBlock(ex)

	var res a.Function
	res = a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		l := ap(c, args)
		l.Put("recur", res)
		return a.Eval(l, db)
	}).WithMetadata(d.meta).(a.Function)
	return res
}

func lambda(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	fn, r := optionalName(args)
	md, r := optionalMetadata(c, r)
	md = loadDocumentation(md)
	an := a.AssertVector(r.First())

	return defineFunction(c, &functionDefinition{
		args: an,
		body: r.Rest(),
		meta: md.Child(a.Properties{
			a.NameKey: fn,
			a.ArgsKey: an,
		}),
	})
}

func apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f := a.AssertApplicable(args.First())
	s := a.AssertSequence(args.Rest().First())
	return f.Apply(c, s)
}

func isApplicable(v a.Value) bool {
	if _, ok := v.(a.Applicable); ok {
		return true
	}
	return false
}

func init() {
	RegisterBuiltIn("lambda", lambda)
	RegisterBuiltIn("apply", apply)
	RegisterSequencePredicate("apply?", isApplicable)
}
