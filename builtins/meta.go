package builtins

import (
	a "github.com/kode4food/sputter/api"
)

// ExpectedName is raised if a Name is expected but not encountered
const ExpectedName = "expected a name"

func toName(v a.Value) a.Name {
	if n, ok := v.(a.Named); ok {
		return n.Name()
	}
	if s, ok := v.(*a.Symbol); ok {
		return s.Name
	}
	if s, ok := v.(string); ok {
		return a.Name(s)
	}
	panic(ExpectedName)
}

func toVariables(args a.Sequence) a.Variables {
	r := make(a.Variables)
	for i := args.(a.Sequence); i.IsSequence(); i = i.Rest() {
		p := a.AssertSequence(i.First())
		a.AssertArity(p, 2)
		k := toName(p.First())
		v := p.Rest().First()
		r[k] = v
	}
	return r
}

func withMeta(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	o := a.AssertAnnotated(a.Eval(c, args.First()))
	m := a.AssertSequence(a.Eval(c, args.Rest().First()))
	return o.WithMetadata(toVariables(m))
}

func init() {
	registerAnnotated(
		a.NewFunction(withMeta).WithMetadata(a.Variables{
			a.MetaName: a.Name("with-meta"),
		}),
	)
}
