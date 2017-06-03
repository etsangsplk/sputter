package evaluator_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	_ "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

func TestCreateExpander(t *testing.T) {
	as := assert.New(t)
	l := e.NewLexer("99")
	c := a.NewContext()
	tr := e.NewReader(c, l)
	ex := e.Expand(c, tr)
	as.NotNil(ex)
}

func TestExpander(t *testing.T) {
	as := assert.New(t)

	b := e.NewEvalContext()
	ns := a.GetContextNamespace(b)
	ns.Delete("hello")
	ns.Put("hello", a.NewMacro(func(_ a.Context, l a.Sequence) a.Value {
		if _, ok := l.(a.List); !ok {
			as.Fail("provided sequence is not a list")
		}
		return a.NewVector(s("you"))
	}).WithMetadata(a.Metadata{
		a.MetaName: a.Name("hello"),
	}).(a.Function))

	l := e.NewLexer(`(hello)`)
	tr := e.NewReader(b, l)
	ex := e.ExpandSequence(b, tr)
	v := a.EvalBlock(b, ex)

	if rv, ok := v.(a.Vector); ok {
		v1, ok := rv.ElementAt(0)
		as.True(ok)
		as.String("you", v1)
	} else {
		as.Fail("expand did not transform")
	}
}

func TestSimpleData(t *testing.T) {
	as := assert.New(t)

	c := e.NewEvalContext()
	v1, ok := e.EvalStr(c, `'99`).(a.Number)
	as.True(ok)
	as.Number(99, v1)
}

func TestSimpleExpand(t *testing.T) {
	as := assert.New(t)

	c := e.NewEvalContext()
	v2, ok := e.EvalStr(c, `[1 2 3]`).(a.Vector)
	as.True(ok)
	as.String("[1 2 3]", v2)

	c = e.NewEvalContext()
	v3, ok := e.EvalStr(c, `{:name "bob"}`).(a.Associative)
	as.True(ok)
	as.String(`{:name "bob"}`, v3)
}

func TestNestedData(t *testing.T) {
	as := assert.New(t)

	c := e.NewEvalContext()
	v := e.EvalStr(c, `''99`)
	as.String("(sputter:quote 99)", v)
}

func TestListData(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	v := e.EvalStr(c, `'(symbol true)`)
	value, ok := v.(a.List)
	as.True(ok)

	if sym, ok := value.First().(a.Symbol); ok {
		as.String("symbol", sym.Name())
	} else {
		as.Fail("first element should be symbol")
	}

	if n, ok := value.Rest().(a.List); ok {
		b, ok := n.First().(a.Value)
		as.True(ok)
		as.Equal(a.True, b)

		nl, ok := n.Rest().(a.List)
		as.True(ok)
		as.Equal(a.EmptyList, nl)
		as.False(nl.IsSequence())
	} else {
		as.Fail("rest() elements not a list")
	}
}
