package main_test

import (
	"testing"

	s "github.com/kode4food/sputter"
	"github.com/stretchr/testify/assert"
)

func assertGet(a *assert.Assertions, c *s.Context, key string, value s.Value) {
	v, ok := c.Get(key)
	a.True(ok)
	a.Equal(value, v)
}

func assertMissing(a *assert.Assertions, c *s.Context, key string) {
	v, ok := c.Get(key)
	a.False(ok)
	a.Equal(s.EmptyList, v)
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

	c2 := c1.Child()
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
	sg2 := sg1.Child()
	sg3 := sg2.Child()

	a.Equal(sg1, sg3.Globals())
}
