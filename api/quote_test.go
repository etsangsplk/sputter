package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestQuote(t *testing.T) {
	as := assert.New(t)

	f := a.NewFloat(99.0)
	q := a.Quote(f)
	c := a.NewContext()

	as.Equal(f, q.Eval(c), "wrapped value returned with Eval()")
	as.Equal("99", a.String(q), "string returned")
}
