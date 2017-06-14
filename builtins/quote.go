package builtins

import (
	"fmt"
	"strings"
	"sync/atomic"

	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
	e "github.com/kode4food/sputter/evaluator"
)

const genSymTemplate = "x-%s-gensym-%d"

var (
	sQuote  = a.NewQualifiedSymbol("quote", a.BuiltInDomain)
	sList   = a.NewQualifiedSymbol("list", a.BuiltInDomain)
	sVector = a.NewQualifiedSymbol("vector", a.BuiltInDomain)
	sAssoc  = a.NewQualifiedSymbol("assoc", a.BuiltInDomain)
)

var genSymIncrement uint64

type syntaxContext struct {
	context a.Context
	genSyms map[string]a.Value
}

func (sc *syntaxContext) quote(v a.Value) a.Value {
	return sc.quoteValue(v)
}

func (sc *syntaxContext) quoteValue(v a.Value) a.Value {
	if s, ok := v.(a.Sequence); ok {
		return sc.quoteSequence(s)
	}
	if s, ok := v.(a.Symbol); ok {
		return sc.quoteSymbol(s)
	}
	return v
}

func (sc *syntaxContext) quoteSymbol(s a.Symbol) a.Value {
	if gs, ok := sc.genSym(s); ok {
		return gs
	}
	return sc.qualifySymbol(s)
}

func (sc *syntaxContext) genSym(s a.Symbol) (a.Value, bool) {
	if s.Domain() != a.LocalDomain {
		return nil, false
	}

	n := string(s.Name())
	if len(n) <= 1 || !strings.HasSuffix(n, "#") {
		return nil, false
	}

	if r, ok := sc.genSyms[n]; ok {
		return r, true
	}

	idx := atomic.AddUint64(&genSymIncrement, 1)
	q := fmt.Sprintf(genSymTemplate, n[0:len(n)-1], idx)
	r := a.NewList(sQuote, a.NewLocalSymbol(a.Name(q)))
	sc.genSyms[n] = r
	return r, true
}

func (sc *syntaxContext) qualifySymbol(s a.Symbol) a.Value {
	if s.Domain() != a.LocalDomain {
		return a.NewList(sQuote, s)
	}
	n := s.Name()
	if c, ok := sc.context.Has(s.Name()); ok {
		if c != sc.context {
			ns := a.GetContextNamespace(c)
			return a.NewList(sQuote, ns.Intern(n))
		}
	}
	return s
}

func (sc *syntaxContext) quoteSequence(s a.Sequence) a.Value {
	var sym a.Symbol
	if l, ok := s.(a.List); ok && l.Count() > 0 {
		sym = sList
	} else if v, ok := s.(a.Vector); ok && v.Count() > 0 {
		sym = sVector
	} else if an, ok := s.(a.Associative); ok && an.Count() > 0 {
		sym = sAssoc
	} else {
		return s
	}
	return a.NewList(sc.quoteElements(s)...).Prepend(sym)
}

func (sc *syntaxContext) quoteElements(s a.Sequence) []a.Value {
	r := []a.Value{}
	for i := s; i.IsSequence(); i = i.Rest() {
		e := sc.quoteValue(i.First())
		r = append(r, e)
	}
	return r
}

func quote(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	return args.First()
}

func unquote(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	return e.Expand(c, args.First())
}

func unquoteSplicing(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	return e.NewSplice(e.Expand(c, args.First()))
}

func syntaxquote(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	sc := &syntaxContext{
		context: c,
		genSyms: make(map[string]a.Value),
	}
	q := sc.quote(args.First())
	r := e.Expand(c, q)
	return r
}

func init() {
	registerAnnotated(
		a.NewMacro(quote).WithMetadata(a.Metadata{
			a.MetaName: a.Name("quote"),
			a.MetaDoc:  d.Get("quote"),
		}),
	)

	registerAnnotated(
		a.NewMacro(unquote).WithMetadata(a.Metadata{
			a.MetaName: a.Name("unquote"),
		}),
	)

	registerAnnotated(
		a.NewMacro(unquoteSplicing).WithMetadata(a.Metadata{
			a.MetaName: a.Name("unquote-splicing"),
		}),
	)

	registerAnnotated(
		a.NewMacro(syntaxquote).WithMetadata(a.Metadata{
			a.MetaName: a.Name("syntax-quote"),
		}),
	)
}
