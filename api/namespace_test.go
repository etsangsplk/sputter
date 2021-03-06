package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestNamespace(t *testing.T) {
	as := assert.New(t)
	ns1 := a.GetNamespace("user")
	ns2 := a.GetNamespace(a.UserDomain)

	as.Equal(ns1, ns2)
	as.Equal(a.UserDomain, ns1.Domain())
	as.Equal(a.Name("user"), ns2.Domain())
	as.String("(ns user)", ns1)
}

func TestWithNamespace(t *testing.T) {
	as := assert.New(t)

	ns1 := a.GetNamespace(a.UserDomain)
	ns2 := a.GetNamespace("foo")

	ns1.Delete("foo")
	ns2.Delete("foo")

	c1 := a.ChildContext(ns1, a.Variables{})
	c2 := a.WithNamespace(c1, ns2)

	ns1.Put("foo", s("outer"))
	c1.Put("bar", s("skipped"))
	ns2.Put("foo", s("inner"))

	v1, _ := c1.Get("foo")
	v2, _ := c2.Get("foo")
	v3, _ := c2.Get("bar")

	as.String("outer", v1)
	as.String("inner", v2)
	as.String("skipped", v3)
}

//func TestWriteOnceRebind(t *testing.T) {
//	as := assert.New(t)
//
//	c := a.NewWriteOnceContext()
//	c.Put("hello", s("there"))
//
//	defer as.ExpectError(a.ErrStr(a.AlreadyBound, a.Name("hello")))
//	c.Put("hello", s("twice"))
//}
