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

func TestAssertNamespace(t *testing.T) {
	as := assert.New(t)
	a.AssertNamespace(a.GetNamespace("hello"))

	defer expectError(as, a.ExpectedNamespace)
	a.AssertNamespace("hello")
}
