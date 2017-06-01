package builtins

import a "github.com/kode4food/sputter/api"

type closure struct {
	names []a.Name
	body  a.Value
}

var emptyNames = []a.Name{}

func makeNames(s a.Sequence) []a.Name {
	v := a.AssertVector(s)
	l := v.Count()
	r := make([]a.Name, l)
	for i := 0; i < l; i++ {
		v, _ := v.ElementAt(i)
		r[i] = a.AssertUnqualified(v).Name()
	}
	return r
}

func consolidateNames(include []a.Name, exclude []a.Name) []a.Name {
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
	r := []a.Name{}
	for k, v := range s {
		if v {
			r = append(r, a.Name(k))
		}
	}
	return r
}

func visitValue(v a.Value) []a.Name {
	if s, ok := v.(a.Sequence); ok {
		return visitSequence(s)
	}
	if _, ok := v.(a.Expression); !ok {
		return emptyNames
	}
	if s, ok := v.(a.Symbol); ok && s.Domain() == a.LocalDomain {
		return []a.Name{s.Name()}
	}
	return emptyNames
}

func visitSequence(s a.Sequence) []a.Name {
	if _, ok := s.(a.Str); ok {
		return emptyNames
	}
	r := []a.Name{}
	for i := s; i.IsSequence(); i = i.Rest() {
		n := visitValue(i.First())
		r = append(r, n...)
	}
	return r
}

func makeClosure(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	ex := makeNames(a.AssertVector(args.First()))
	body := evalExpandBlock(c, args.Rest())
	names := consolidateNames(visitValue(body), ex)

	if len(names) > 0 {
		return &closure{
			names: names,
			body:  body,
		}
	}
	return body
}

func (cl *closure) Eval(c a.Context) a.Value {
	names := cl.names
	vars := make(a.Variables, len(names))
	for _, n := range names {
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
