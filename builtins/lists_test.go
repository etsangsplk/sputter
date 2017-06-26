package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestList(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("x")

	testCode(t, `(list? '(1 2 3))`, a.True)
	testCode(t, `(list? ())`, a.True)
	testCode(t, `(list? [1 2 3])`, a.False)
	testCode(t, `(list? 42)`, a.False)
	testCode(t, `(list? (list 1 2 3))`, a.True)
	testCode(t, `(list)`, a.EmptyList)

	testCode(t, `
		(def x '(1 2 3 4))
		(x 2)
	`, f(3))
}
