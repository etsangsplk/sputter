package api_test

import (
	"fmt"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func assertGet(as *assert.Assertions, c a.Context, n a.Name, cv a.Value) {
	v, ok := c.Get(n)
	as.True(ok)
	as.Equal(cv, v)
}

func assertMissing(as *assert.Assertions, c a.Context, n a.Name) {
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
	c.Put("hello", "there")
	assertGet(as, c, "hello", "there")
}

func TestNestedContext(t *testing.T) {
	as := assert.New(t)

	c1 := a.NewContext()
	c1.Put("hello", "there")
	c1.Put("howdy", "ho")

	c2 := a.ChildContext(c1)
	c2.Put("hello", "you")
	c2.Put("foo", "bar")

	assertGet(as, c1, "hello", "there")
	assertGet(as, c1, "howdy", "ho")
	assertMissing(as, c1, "foo")

	assertGet(as, c2, "hello", "you")
	assertGet(as, c2, "howdy", "ho")
	assertGet(as, c2, "foo", "bar")
}

func TestEvalContext(t *testing.T) {
	as := assert.New(t)

	uc := a.GetNamespace(a.UserDomain)
	uc.Delete("foo")
	uc.Put("foo", 99)

	ec := a.NewEvalContext()
	v, _ := ec.Get("foo")
	as.Equal(99, v, "EvalContext chain works")
}

func TestRebind(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	c.Put("hello", "there")

	defer expectError(as, fmt.Sprintf(a.AlreadyBound, "hello"))
	c.Put("hello", "twice")
}
