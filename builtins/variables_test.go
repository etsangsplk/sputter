package builtins_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestVariables(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("foo")
	ns.Delete("return-local")

	testCode(t, `
		(def foo "bar")
		foo
	`, "bar")

	testCode(t, `
		(defn return-local []
			(let [foo "local"] foo))
		(return-local)
	`, "local")
}

func TestScopeQualifiers(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("foo")
	testCode(t, `
		(def foo 99)
		(let [foo 100]
			(+ foo user:foo))
	`, big.NewFloat(199))
}
