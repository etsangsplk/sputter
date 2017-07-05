package builtins

import (
	"fmt"
	"strings"
	"sync/atomic"

	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

const (
	// UnsupportedSyntaxQuote is raised when something can't be quoted
	UnsupportedSyntaxQuote = "unsupported type in syntax quote: %s"

	genSymTemplate = "x-%s-gensym-%d"
)

type syntaxContext struct {
	context a.Context
	genSyms map[string]a.Symbol
}

var (
	sQuote  = a.NewQualifiedSymbol("quote", a.BuiltInDomain)
	sList   = a.NewQualifiedSymbol("list", a.BuiltInDomain)
	sVector = a.NewQualifiedSymbol("vector", a.BuiltInDomain)
	sAssoc  = a.NewQualifiedSymbol("assoc", a.BuiltInDomain)
	sApply  = a.NewQualifiedSymbol("apply", a.BuiltInDomain)
	sAppend = a.NewQualifiedSymbol("append", a.BuiltInDomain)

	genSymIncrement uint64
)

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
	if gs, ok := sc.generateSymbol(s); ok {
		return a.NewList(sQuote, gs)
	}
	return a.NewList(sQuote, sc.qualifySymbol(s))
}

func (sc *syntaxContext) generateSymbol(s a.Symbol) (a.Symbol, bool) {
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
	r := a.NewLocalSymbol(a.Name(q))
	sc.genSyms[n] = r
	return r, true
}

func (sc *syntaxContext) qualifySymbol(s a.Symbol) a.Value {
	if s.Domain() != a.LocalDomain {
		return s
	}
	n := s.Name()
	if c, ok := sc.context.Has(s.Name()); ok {
		if c != sc.context {
			ns := a.GetContextNamespace(c)
			return ns.Intern(n)
		}
	}
	return s
}

func (sc *syntaxContext) quoteSequence(s a.Sequence) a.Value {
	if st, ok := s.(a.Str); ok {
		return st
	}
	if l, ok := s.(a.List); ok {
		return a.NewList(sApply, sList, sc.quoteElements(l))
	}
	if v, ok := s.(a.Vector); ok {
		return a.NewList(sApply, sVector, sc.quoteElements(v))
	}
	if as, ok := s.(a.Associative); ok {
		return sc.quoteAssociative(as)
	}
	panic(a.Err(UnsupportedSyntaxQuote, s))
}

func (sc *syntaxContext) quoteAssociative(as a.Associative) a.Value {
	r := []a.Value{}
	for i := as.(a.Sequence); i.IsSequence(); i = i.Rest() {
		p := i.First().(a.Vector)
		k, _ := p.ElementAt(0)
		v, _ := p.ElementAt(1)
		r = append(r, k)
		r = append(r, v)
	}
	return a.NewList(sApply, sAssoc, sc.quoteElements(a.NewVector(r...)))
}

func (sc *syntaxContext) quoteElements(s a.Sequence) a.Value {
	r := []a.Value{}
	for i := s; i.IsSequence(); i = i.Rest() {
		v := i.First()
		if f, ok := isUnquoteSplicing(v); ok {
			r = append(r, f)
			continue
		}
		if f, ok := isUnquote(v); ok {
			r = append(r, a.NewList(sList, f))
			continue
		}
		r = append(r, a.NewList(sList, sc.quoteValue(v)))
	}
	return a.NewList(r...).Prepend(sAppend)
}

func isWrapperCall(n a.Name, v a.Value) (a.Value, bool) {
	if l, ok := isBuiltInCall(n, v); ok {
		return l.Rest().First(), true
	}
	return nil, false
}

func isUnquote(v a.Value) (a.Value, bool) {
	return isWrapperCall("unquote", v)
}

func isUnquoteSplicing(v a.Value) (a.Value, bool) {
	return isWrapperCall("unquote-splicing", v)
}

func quote(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	return args.First()
}

func syntaxquote(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	sc := &syntaxContext{
		context: c,
		genSyms: make(map[string]a.Symbol),
	}
	return sc.quote(args.First())
}

func init() {
	registerAnnotated(
		a.NewMacro(quote).WithMetadata(a.Metadata{
			a.MetaName:    a.Name("quote"),
			a.MetaDoc:     d.Get("quote"),
			a.MetaSpecial: a.True,
		}),
	)

	registerAnnotated(
		a.NewMacro(syntaxquote).WithMetadata(a.Metadata{
			a.MetaName: a.Name("syntax-quote"),
		}),
	)
}
