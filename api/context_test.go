package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func assertGet(as *assert.Wrapper, c a.Context, n a.Name, cv a.Value) {
	v, ok := c.Get(n)
	as.True(ok)
	as.Equal(cv, v)
}

func assertMissing(as *assert.Wrapper, c a.Context, n a.Name) {
	v, ok := c.Get(n)
	as.False(ok)
	as.Equal(a.Nil, v)
}

func TestCreateContext(t *testing.T) {
	as := assert.New(t)
	c := a.NewContext()
	as.NotNil(c)
}

func TestPopulateContext(t *testing.T) {
	as := assert.New(t)
	c := a.NewContext()
	c.Put("hello", s("there"))
	assertGet(as, c, "hello", s("there"))
}

func TestNestedContext(t *testing.T) {
	as := assert.New(t)

	c1 := a.NewContext()
	c1.Put("hello", s("there"))
	c1.Put("howdy", s("ho"))

	c2 := a.ChildContext(c1)
	c2.Put("hello", s("you"))
	c2.Put("foo", s("bar"))

	assertGet(as, c1, "hello", s("there"))
	assertGet(as, c1, "howdy", s("ho"))
	assertMissing(as, c1, "foo")

	assertGet(as, c2, "hello", s("you"))
	assertGet(as, c2, "howdy", s("ho"))
	assertGet(as, c2, "foo", s("bar"))
}

func TestRebind(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	c.Put("hello", s("there"))

	defer as.ExpectError(a.Err(a.AlreadyBound, a.Name("hello")))
	c.Put("hello", s("twice"))
}
