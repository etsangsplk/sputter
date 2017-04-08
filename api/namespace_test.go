package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestNamespace(t *testing.T) {
	as := assert.New(t)
	ns1 := a.GetNamespace("user")
	ns2 := a.GetNamespace(a.UserDomain)

	as.Equal(ns1, ns2, "Namespaces identical")
	as.Equal(a.UserDomain, ns1.Domain(), "correct domain")
	as.Equal(a.Name("user"), ns2.Domain(), "correct domain")
	as.Equal("(ns user)", a.String(ns1), "correct string representation")
}

func TestWithNamespace(t *testing.T) {
	as := assert.New(t)

	ns1 := a.GetNamespace(a.UserDomain)
	ns2 := a.GetNamespace("foo")

	ns1.Delete("foo")
	ns2.Delete("foo")

	ns1.Put("foo", "outer")
	c1 := a.ChildContext(ns1)
	c2 := a.WithNamespace(c1, ns2)
	ns2.Put("foo", "inner")

	v1, _ := c1.Get("foo")
	v2, _ := c2.Get("foo")

	as.Equal("outer", v1, "outer is correct")
	as.Equal("inner", v2, "inner is correct")
}

func TestAssertNamespace(t *testing.T) {
	as := assert.New(t)
	a.AssertNamespace(a.GetNamespace("hello"))

	defer expectError(as, a.Err(a.ExpectedNamespace, a.String("hello")))
	a.AssertNamespace("hello")
}
