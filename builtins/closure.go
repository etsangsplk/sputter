package builtins

import a "github.com/kode4food/sputter/api"

const (
	makeClosureName = "make-closure"
	closureName     = "closure"
)

type (
	makeClosureFunction struct{ BaseBuiltIn }
	closureFunction     struct{ BaseBuiltIn }
)

var (
	closureSym = a.NewBuiltInSymbol(closureName)
	emptyNames = a.Names{}
)

func assertUnqualifiedNames(s a.Sequence) a.Names {
	v := s.(a.Vector)
	l := v.Count()
	r := make(a.Names, l)
	for i := 0; i < l; i++ {
		e, _ := v.ElementAt(i)
		r[i] = e.(a.LocalSymbol).Name()
	}
	return r
}

func makeLocalSymbolVector(names a.Names) a.Vector {
	nl := len(names)
	nv := make(a.Vector, nl)
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

func (*makeClosureFunction) Apply(c a.Context, args a.Vector) a.Value {
	a.AssertMinimumArity(args, 1)
	ex := assertUnqualifiedNames(args[0].(a.Vector))
	cb := args[1:]
	nm := consolidateNames(visitValue(cb), ex)
	return a.NewList(closureSym, makeLocalSymbolVector(nm), cb)
}

func isClosure(v a.Value) (a.Names, bool) {
	if l, ok := isBuiltInCall(closureName, v); ok {
		e := l.Rest().First().(a.Vector)
		return assertUnqualifiedNames(e), true
	}
	return emptyNames, false
}

func (*closureFunction) Apply(c a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 2)
	v := args[0].(a.Vector)
	s := args[1].(a.Sequence)
	bl := a.MakeBlock(s)
	return closureFromVector(c, v, bl)
}

func closureFromVector(c a.Context, v a.Vector, bl a.Block) a.Value {
	names := make(a.Names, len(v))
	for i, n := range v {
		names[i] = n.(a.LocalSymbol).Name()
	}
	l := a.MakeClosure(c, names)
	return bl.Eval(l)
}

func init() {
	var makeClosure *makeClosureFunction
	var closure *closureFunction

	RegisterBuiltIn(makeClosureName, makeClosure)
	RegisterBuiltIn(closureName, closure)
}
