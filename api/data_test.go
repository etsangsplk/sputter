package api_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestAtom(t *testing.T) {
	a := assert.New(t)

	v := &s.Atom{Label: "hello"}
	a.Equal(v, v.Eval(s.NewContext()), "own value returned")
	a.Equal("hello", v.String(), "label returned")
}

func TestQuote(t *testing.T) {
	a := assert.New(t)

	f := big.NewFloat(99.0)
	q := &s.Quote{Value: f}
	a.Equal(f, q.Eval(s.NewContext()), "wrapped value returned")
	a.Equal("99", q.String(), "string returned")
}
