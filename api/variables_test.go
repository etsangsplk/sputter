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
	as.True(a.True.Apply(c, a.Vector{a.True, a.True, a.True}))
	as.False(a.True.Apply(c, a.Vector{a.True, a.False, a.True}))
	as.True(a.False.Apply(c, a.Vector{a.False, a.False, a.False}))
	as.False(a.False.Apply(c, a.Vector{a.True, a.False, a.True}))
	as.True(a.AssertBool(a.True))

	defer as.ExpectError(a.Err(a.ExpectedBool, s("not bool")))
	a.AssertBool(s("not bool"))
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
	as.True(a.Nil.Apply(c, a.Vector{a.Nil, a.Nil, a.Nil}))
	as.False(a.Nil.Apply(c, a.Vector{a.Nil, a.False, a.Nil}))
}
