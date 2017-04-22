package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) *a.Number {
	return a.NewFloat(f)
}

func TestTruthy(t *testing.T) {
	as := assert.New(t)

	as.Truthy(a.True)
	as.Truthy(a.NewList(a.Str("Hello")))
	as.Truthy(a.Str("hello"))

	as.Falsey(a.Nil)
	as.Falsey(a.False)
}
