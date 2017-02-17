package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
	"math/big"
)

func TestData(t *testing.T) {
	a := assert.New(t)

	f := big.NewFloat(99.0)
	d := &s.Data{Value: f}
	a.Equal(f, d.Eval(s.NewContext()), "wrapped value returned")
	a.Equal("99", d.String(), "string returned")
}
