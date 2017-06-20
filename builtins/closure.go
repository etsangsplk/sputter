package builtins

import (
	a "github.com/kode4food/sputter/api"
	e "github.com/kode4food/sputter/evaluator"
)

type names []a.Name

type closure struct {
	names names
	body  a.Value
}

var emptyNames = names{}

func makeNames(s a.Sequence) names {
	v := a.AssertVector(s)
	l := v.Count()
	r := make(names, l)
	for i := 0; i < l; i++ {
		v, _ := v.ElementAt(i)
		r[i] = a.AssertUnqualified(v).Name()
	}
	return r
}

func consolidateNames(include names, exclude names) names {
	s := map[string]bool{}
	for _, n := range exclude {
		s[string(n)] = false
	}
	for _, n := range include {
		sn := string(n)
		if _, ok := s[sn]; !ok {
			s[sn] = true
		}
	}
	r := names{}
	for k, v := range s {
		if v {
			r = append(r, a.Name(k))
		}
	}
	return r
}

func visitValue(v a.Value) names {
	if s, ok := v.(a.Sequence); ok {
		return visitSequence(s)
	}
	if cl, ok := v.(*closure); ok {
		return cl.names
	}
	if s, ok := v.(a.Symbol); ok && s.Domain() == a.LocalDomain {
		return names{s.Name()}
	}
	return emptyNames
}

func visitSequence(s a.Sequence) names {
	if _, ok := s.(a.Str); ok {
		return emptyNames
	}
	r := names{}
	for i := s; i.IsSequence(); i = i.Rest() {
		n := visitValue(i.First())
		r = append(r, n...)
	}
	return r
}

func makeClosure(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	nv := makeNames(a.AssertVector(args.First()))
	cb := e.EvalExpand(c, args.Rest())
	cn := consolidateNames(visitValue(cb), nv)

	if len(cn) > 0 {
		return &closure{
			names: cn,
			body:  cb,
		}
	}
	return cb
}

func (cl *closure) Eval(c a.Context) a.Value {
	cn := cl.names
	vars := make(a.Variables, len(cn))
	for _, n := range cn {
		if v, ok := c.Get(n); ok {
			vars[n] = v
		}
	}

	ns := a.GetContextNamespace(c)
	l := a.ChildContextVars(ns, vars)
	return cl.body.Eval(l)
}

func (cl *closure) Str() a.Str {
	return a.MakeDumpStr(cl)
}

func init() {
	registerAnnotated(
		a.NewMacro(makeClosure).WithMetadata(a.Metadata{
			a.MetaName: a.Name("closure"),
		}),
	)
}
