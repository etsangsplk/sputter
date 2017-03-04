package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestSymbol(t *testing.T) {
	a := assert.New(t)

	c := s.NewContext()
	c.Put("howdy", "ho")

	sym := &s.Symbol{Name: "howdy"}
	a.Equal("ho", sym.Eval(c), "symbol value returned")
	a.Equal("howdy", sym.String(), "symbol name returned")
}

func TestQualifiedSymbol(t *testing.T) {
	a := assert.New(t)

	ns1 := s.GetNamespace("ns1")
	ns1.Put("foo", "foo-ns1")
	ns2 := s.GetNamespace("ns2")
	ns2.Put("foo", "foo-ns2")

	empty := s.NewContext()

	c1 := s.NewContext()
	c1.Put(s.ContextDomain, s.Name("ns1"))

	c2 := s.NewContext()
	c2.Put(s.ContextDomain, s.Name("ns2"))

	s1 := s.ParseSymbol("ns1:foo")
	s2 := s.ParseSymbol("ns2:foo")
	s3 := s.ParseSymbol("foo")

	a.Equal("foo-ns1", s1.Eval(empty))
	a.Equal("foo-ns2", s2.Eval(empty))
	a.Equal("foo-ns1", s3.Eval(c1))
	a.Equal("foo-ns2", s3.Eval(c2))
}

func TestSymbolInterning(t *testing.T) {
	a := assert.New(t)

	sym1 := s.NewLocalSymbol("hello")
	sym2 := s.NewLocalSymbol("there")
	sym3 := s.NewLocalSymbol("hello")

	a.Equal(sym1, sym3, "properly interned")
	a.NotEqual(sym1, sym2, "properly isolated")
}

func TestUnknownSymbol(t *testing.T) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(s.UnknownSymbol, rec, "symbol eval panics properly")
			return
		}
		a.Fail("symbol eval should panic")
	}()

	c := s.NewContext()
	sym := &s.Symbol{Name: "howdy"}
	sym.Eval(c)
}

func TestSymbolParsing(t *testing.T) {
	a := assert.New(t)

	s1 := s.ParseSymbol("domain:name1")
	a.Equal("domain", string(s1.Domain))
	a.Equal("name1", string(s1.Name))

	s2 := s.ParseSymbol(":name2")
	a.Equal(s.LocalDomain, s2.Domain)
	a.Equal("name2", string(s2.Name))

	s3 := s.ParseSymbol("name3")
	a.Equal(s.LocalDomain, s3.Domain)
	a.Equal("name3", string(s3.Name))

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(s.BadQualifiedName, rec, "parse explodes correctly")
			return
		}
		a.Fail("bad parse should explode")
	}()

	s.ParseSymbol("one:too:")
}
