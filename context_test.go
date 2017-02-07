package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertGet(a *assert.Assertions, c *Context, key string, value Value) {
	v, ok := c.Get(key)
	a.True(ok)
	a.Equal(value, v)
}

func assertMissing(a *assert.Assertions, c *Context, key string) {
	v, ok := c.Get(key)
	a.False(ok)
	a.Equal(EmptyList, v)
}

func TestCreateContext(t *testing.T) {
	a := assert.New(t)
	c := NewContext()
	a.NotNil(c)
}

func TestPopulateContext(t *testing.T) {
	a := assert.New(t)
	c := NewContext()
	c.Put("hello", "there")
	assertGet(a, c, "hello", "there")
}

func TestNestedContext(t *testing.T) {
	a := assert.New(t)

	c1 := NewContext()
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

	sg1 := NewContext()
	sg2 := sg1.Child()
	sg3 := sg2.Child()

	a.Equal(sg1, sg3.Globals())
}
