package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestNamespace(t *testing.T) {
	a := assert.New(t)
	ns1 := s.GetNamespace("user")
	ns2 := s.GetNamespace(s.UserDomain)

	a.Equal(ns1, ns2, "Namespaces identical")
	a.Equal(s.UserDomain, ns1.Domain(), "correct domain")
	a.Equal(s.Name("user"), ns2.Domain(), "correct domain")
}
