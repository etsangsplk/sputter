package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestSymbol(t *testing.T) {
	as := assert.New(t)

	c := a.Variables{}
	c.Put("howdy", s("ho"))

	sym := a.NewLocalSymbol("howdy")
	as.String("ho", a.Eval(c, sym))
	as.String("howdy", sym)
}

func TestQualifiedSymbol(t *testing.T) {
	as := assert.New(t)

	ns1 := a.GetNamespace("ns1")
	ns1.Put("foo", s("foo-ns1"))
	ns2 := a.GetNamespace("ns2")
	ns2.Put("foo", s("foo-ns2"))

	empty := a.Variables{}

	c1 := a.Variables{}
	c1.Put(a.ContextDomain, a.GetNamespace("ns1"))

	c2 := a.Variables{}
	c2.Put(a.ContextDomain, a.GetNamespace("ns2"))

	s1 := a.ParseSymbol("ns1:foo")
	s2 := a.ParseSymbol("ns2:foo")
	s3 := a.ParseSymbol("foo")

	as.Equal(a.GetNamespace("ns1"), s1.Namespace(c2))
	as.Equal(a.GetNamespace("ns2"), s2.Namespace(c1))
	as.Equal(a.GetNamespace("ns1"), s3.Namespace(c1))

	as.String("foo-ns1", a.Eval(empty, s1))
	as.String("foo-ns2", a.Eval(empty, s2))
	as.String("foo-ns1", a.Eval(c1, s3))
	as.String("foo-ns2", a.Eval(c2, s3))
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

	defer as.ExpectError(a.ErrStr(a.UnknownSymbol, "howdy"))
	c := a.Variables{}
	sym := a.NewLocalSymbol("howdy")
	a.Eval(c, sym)
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
