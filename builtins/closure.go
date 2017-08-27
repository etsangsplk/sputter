package builtins

import a "github.com/kode4food/sputter/api"

type (
	makeClosureFunction struct{ BaseBuiltIn }
	closureFunction     struct{ BaseBuiltIn }
)

var (
	closureSym = a.NewBuiltInSymbol("closure")
	emptyNames = a.Names{}
)

func assertUnqualifiedNames(s a.Sequence) a.Names {
	v := a.AssertVector(s)
	l := v.Count()
	r := make(a.Names, l)
	for i := 0; i < l; i++ {
		v, _ := v.ElementAt(i)
		r[i] = a.AssertUnqualified(v).Name()
	}
	return r
}

func makeLocalSymbolVector(names a.Names) a.Vector {
	nl := len(names)
	nv := make(a.Values, nl)
	for i := 0; i < nl; i++ {
		nv[i] = a.NewLocalSymbol(names[i])
	}
	return nv
}

func consolidateNames(include a.Names, exclude a.Names) a.Names {
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
	r := a.Names{}
	for k, v := range s {
		if v {
			r = append(r, a.Name(k))
		}
	}
	return r
}

func visitValue(v a.Value) a.Names {
	if s, ok := v.(a.Sequence); ok {
		return visitSequence(s)
	}
	if n, ok := isClosure(v); ok {
		return n
	}
	if s, ok := v.(a.Symbol); ok && s.Domain() == a.LocalDomain {
		return a.Names{s.Name()}
	}
	return emptyNames
}

func visitSequence(s a.Sequence) a.Names {
	if _, ok := s.(a.Str); ok {
		return emptyNames
	}
	res := a.Names{}
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		n := visitValue(f)
		res = append(res, n...)
	}
	return res
}

func (*makeClosureFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	f, r, _ := args.Split()
	ex := assertUnqualifiedNames(a.AssertVector(f))
	cb := a.MacroExpandAll(c, r)
	nm := consolidateNames(visitValue(cb), ex)
	return a.NewList(closureSym, makeLocalSymbolVector(nm), cb)
}

func isClosure(v a.Value) (a.Names, bool) {
	if l, ok := isBuiltInCall("closure", v); ok {
		v := a.AssertVector(l.Rest().First())
		return assertUnqualifiedNames(v), true
	}
	return emptyNames, false
}

func (*closureFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f, r, _ := args.Split()
	in := a.AssertVector(f)
	vars := make(a.Variables, in.Count())
	for f, r, ok := in.(a.Sequence).Split(); ok; f, r, ok = r.Split() {
		n := a.AssertUnqualified(f).Name()
		if v, ok := c.Get(n); ok {
			vars[n] = v
		}
	}

	s := a.AssertSequence(r.First())
	bl := a.MakeBlock(s)
	ns := a.GetContextNamespace(c)
	l := a.ChildContextVars(ns, vars)
	return bl.Eval(l)
}

func init() {
	var makeClosure *makeClosureFunction
	var closure *closureFunction

	RegisterBuiltIn("make-closure", makeClosure)
	RegisterBuiltIn("closure", closure)
}
