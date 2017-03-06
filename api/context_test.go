package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func assertGet(a *assert.Assertions, c s.Context, n s.Name, cv s.Value) {
	v, ok := c.Get(n)
	a.True(ok)
	a.Equal(cv, v)
}

func assertMissing(a *assert.Assertions, c s.Context, n s.Name) {
	v, ok := c.Get(n)
	a.False(ok)
	a.Equal(s.Nil, v)
}

func TestCreateContext(t *testing.T) {
	a := assert.New(t)
	c := s.NewContext()
	a.NotNil(c)
}

func TestPopulateContext(t *testing.T) {
	a := assert.New(t)
	c := s.NewContext()
	c.Put("hello", "there")
	assertGet(a, c, "hello", "there")
}

func TestNestedContext(t *testing.T) {
	a := assert.New(t)

	c1 := s.NewContext()
	c1.Put("hello", "there")
	c1.Put("howdy", "ho")

	c2 := s.ChildContext(c1)
	c2.Put("hello", "you")
	c2.Put("foo", "bar")

	assertGet(a, c1, "hello", "there")
	assertGet(a, c1, "howdy", "ho")
	assertMissing(a, c1, "foo")

	assertGet(a, c2, "hello", "you")
	assertGet(a, c2, "howdy", "ho")
	assertGet(a, c2, "foo", "bar")
}

func TestEvalContext(t *testing.T) {
	a := assert.New(t)

	uc := s.GetNamespace(s.UserDomain)
	uc.Delete("foo")
	uc.Put("foo", 99)

	ec := s.NewEvalContext()
	v, _ := ec.Get("foo")
	a.Equal(99, v, "EvalContext chain works")
}
