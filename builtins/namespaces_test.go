package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestNamespaces(t *testing.T) {
	ns1 := a.GetNamespace("foo")
	ns2 := a.GetNamespace("bar")

	ns1.Delete("v1")
	ns2.Delete("v1")

	testCode(t, `
		(with-ns foo
			(def v1 99))
		(with-ns bar
			(def v1 100)
			(+ v1 foo:v1))
	`, a.NewFloat(199))

	testCode(t, `(ns foo)`, ns1)

	testBadCode(t, `
		(ns foo:bar)
	`, a.ExpectedUnqualified)
}
