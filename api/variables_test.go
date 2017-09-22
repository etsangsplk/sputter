package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestNames(t *testing.T) {
	as := assert.New(t)

	n := a.Name("hello")
	as.Equal(a.Name("hello"), n.Name())
}

func TestBool(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	as.True(a.True.Apply(c, a.NewVector(a.True, a.True, a.True)))
	as.False(a.True.Apply(c, a.NewVector(a.True, a.False, a.True)))
	as.True(a.False.Apply(c, a.NewVector(a.False, a.False, a.False)))
	as.False(a.False.Apply(c, a.NewVector(a.True, a.False, a.True)))
}

func TestTruthy(t *testing.T) {
	as := assert.New(t)

	as.Truthy(a.True)
	as.Truthy(a.NewList(a.Str("Hello")))
	as.Truthy(a.Str("hello"))

	as.Falsey(a.Nil)
	as.Falsey(a.False)
}

func TestNil(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	as.True(a.Nil.Apply(c, a.NewVector(a.Nil, a.Nil, a.Nil)))
	as.False(a.Nil.Apply(c, a.NewVector(a.Nil, a.False, a.Nil)))
}
