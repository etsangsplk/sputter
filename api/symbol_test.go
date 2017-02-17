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
