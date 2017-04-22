package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestQuote(t *testing.T) {
	as := assert.New(t)

	f := f(99.0)
	q := a.Quote(f)
	c := a.NewContext()

	as.Equal(f, q.Eval(c))
	as.String("99", q)
}

func TestQuoteApply(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	v := a.Vector{f(1), f(2), f(3)}
	q := a.Quote(v).(a.Applicable)
	r := q.Apply(c, a.NewList(f(1)))
	as.Float(2, r)
}
