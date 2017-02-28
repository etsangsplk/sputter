package builtins_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
)

func TestVariables(t *testing.T) {
	s.GetNamespace(s.UserDomain).Delete("foo")
	testCode(t, `
		(defvar foo "bar")
		foo
	`, "bar")

	s.GetNamespace(s.UserDomain).Delete("return-local")
	testCode(t, `
		(defun return-local []
			(let [foo "local"] foo))
		(return-local)
	`, "local")
}

func TestScopeQualifiers(t *testing.T) {
	s.GetNamespace(s.UserDomain).Delete("foo")
	testCode(t, `
		(defvar foo 99)
		(let [foo 100]
			user:foo)
	`, big.NewFloat(99))
}
