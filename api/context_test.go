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

func TestGlobalContext(t *testing.T) {
	a := assert.New(t)

	sg1 := s.NewContext()
	sg2 := s.ChildContext(sg1)
	sg3 := s.ChildContext(sg2)

	a.Equal(sg1, sg3.Globals())
}
