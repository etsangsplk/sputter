package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestSymbol(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	c.Put("howdy", s("ho"))

	sym := a.NewLocalSymbol("howdy")
	as.True(sym.IsSymbol())
	as.String("ho", sym.Expression().Eval(c))
	as.String("howdy", sym)
}

func TestQualifiedSymbol(t *testing.T) {
	as := assert.New(t)

	ns1 := a.GetNamespace("ns1")
	ns1.Put("foo", s("foo-ns1"))
	ns2 := a.GetNamespace("ns2")
	ns2.Put("foo", s("foo-ns2"))

	empty := a.NewContext()

	c1 := a.NewContext()
	c1.Put(a.ContextDomain, a.GetNamespace("ns1"))

	c2 := a.NewContext()
	c2.Put(a.ContextDomain, a.GetNamespace("ns2"))

	s1 := a.ParseSymbol("ns1:foo").Expression().(a.Symbol)
	s2 := a.ParseSymbol("ns2:foo").Expression().(a.Symbol)
	s3 := a.ParseSymbol("foo").Expression().(a.Symbol)

	as.Equal(a.GetNamespace("ns1"), s1.Namespace(c2))
	as.Equal(a.GetNamespace("ns2"), s2.Namespace(c1))
	as.Equal(a.GetNamespace("ns1"), s3.Namespace(c1))

	as.String("foo-ns1", s1.Eval(empty))
	as.String("foo-ns2", s2.Eval(empty))
	as.String("foo-ns1", s3.Eval(c1))
	as.String("foo-ns2", s3.Eval(c2))
}

func TestSymbolInterning(t *testing.T) {
	as := assert.New(t)

	sym1 := a.NewLocalSymbol("hello")
	sym2 := a.NewLocalSymbol("there")
	sym3 := a.NewLocalSymbol("hello")

	as.Identical(sym1, sym3)
	as.NotIdentical(sym1, sym2)
}

func TestUnknownSymbol(t *testing.T) {
	as := assert.New(t)

	defer as.ExpectError(a.Err(a.UnknownSymbol, "howdy"))
	c := a.NewContext()
	sym := a.NewLocalSymbol("howdy").Expression()
	sym.Eval(c)
}

func TestSymbolParsing(t *testing.T) {
	as := assert.New(t)

	s1 := a.ParseSymbol("domain:name1")
	as.String("domain", string(s1.Domain()))
	as.String("name1", string(s1.Name()))

	s2 := a.ParseSymbol(":name2")
	as.Equal(a.LocalDomain, s2.Domain())
	as.String("name2", string(s2.Name()))

	s3 := a.ParseSymbol("name3")
	as.Equal(a.LocalDomain, s3.Domain())
	as.String("name3", string(s3.Name()))

	s4 := a.ParseSymbol("one:too:")
	as.String("one", string(s4.Domain()))
	as.String("too:", string(s4.Name()))

}

func TestAssertSymbol(t *testing.T) {
	as := assert.New(t)
	defer as.ExpectError(a.Err(a.ExpectedSymbol, "37"))
	a.AssertUnqualified(f(37))
}

func TestAssertUnqualified(t *testing.T) {
	as := assert.New(t)
	a.AssertUnqualified(a.NewLocalSymbol("hello"))

	defer as.ExpectError(a.Err(a.ExpectedUnqualified, "bar:hello"))
	a.AssertUnqualified(a.NewQualifiedSymbol("hello", "bar"))
}
