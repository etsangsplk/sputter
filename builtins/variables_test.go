package builtins_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
)

func TestVariables(t *testing.T) {
	s.GetNamespace(s.UserDomain).Delete("foo")
	testCode(t, `
		(def foo "bar")
		foo
	`, "bar")

	s.GetNamespace(s.UserDomain).Delete("return-local")
	testCode(t, `
		(defn return-local []
			(let [foo "local"] foo))
		(return-local)
	`, "local")
}

func TestScopeQualifiers(t *testing.T) {
	s.GetNamespace(s.UserDomain).Delete("foo")
	testCode(t, `
		(def foo 99)
		(let [foo 100]
			(+ foo user:foo))
	`, big.NewFloat(199))
}
