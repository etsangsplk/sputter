package builtins

import (
	"fmt"
	"strings"
	"sync/atomic"

	a "github.com/kode4food/sputter/api"
)

const (
	// UnsupportedSyntaxQuote is raised when something can't be quoted
	UnsupportedSyntaxQuote = "unsupported type in syntax quote: %s"

	genSymTemplate = "x-%s-gensym-%d"
)

type (
	quoteFunction       struct{ BaseBuiltIn }
	syntaxQuoteFunction struct{ BaseBuiltIn }

	syntaxContext struct {
		context a.Context
		genSyms map[string]a.Symbol
	}
)

var (
	quoteSym  = a.NewBuiltInSymbol("quote")
	listSym   = a.NewBuiltInSymbol("list")
	vectorSym = a.NewBuiltInSymbol("vector")
	assocSym  = a.NewBuiltInSymbol("assoc")
	applySym  = a.NewBuiltInSymbol("apply")
	concatSym = a.NewBuiltInSymbol("concat")

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
		return a.NewList(quoteSym, gs)
	}
	return a.NewList(quoteSym, sc.qualifySymbol(s))
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
		return a.NewList(applySym, listSym, sc.quoteElements(l))
	}
	if v, ok := s.(a.Vector); ok {
		return a.NewList(applySym, vectorSym, sc.quoteElements(v))
	}
	if as, ok := s.(a.Associative); ok {
		return sc.quoteAssociative(as)
	}
	panic(a.ErrStr(UnsupportedSyntaxQuote, s))
}

func (sc *syntaxContext) quoteAssociative(as a.Associative) a.Value {
	res := []a.Value{}
	for f, r, ok := as.Split(); ok; f, r, ok = r.Split() {
		p := f.(a.Vector)
		k, _ := p.ElementAt(0)
		v, _ := p.ElementAt(1)
		res = append(res, k)
		res = append(res, v)
	}
	return a.NewList(applySym, assocSym, sc.quoteElements(a.NewVector(res...)))
}

func (sc *syntaxContext) quoteElements(s a.Sequence) a.Value {
	res := []a.Value{}
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if v, ok := isUnquoteSplicing(f); ok {
			res = append(res, v)
			continue
		}
		if v, ok := isUnquote(f); ok {
			res = append(res, a.NewList(listSym, v))
			continue
		}
		res = append(res, a.NewList(listSym, sc.quoteValue(f)))
	}
	return a.NewList(res...).Prepend(concatSym)
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

func (*quoteFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	return args.First()
}

func (*syntaxQuoteFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	sc := &syntaxContext{
		context: c,
		genSyms: make(map[string]a.Symbol),
	}
	return sc.quote(args.First())
}

func init() {
	var quote *quoteFunction
	var syntaxQuote *syntaxQuoteFunction

	RegisterBuiltIn("quote", quote)
	RegisterBuiltIn("syntax-quote", syntaxQuote)
}
