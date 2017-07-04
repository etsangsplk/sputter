package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestChannel(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("c")

	testCode(t, `
		(def c (channel))
		(do-async
			(apply (:emit c) '("hello"))
			(apply (:close c) ()))
		(first (:seq c))
	`, s("hello"))
}

func TestPromise(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("p1")
	ns.Delete("p2")

	testCode(t, `
		(def p1 (promise))
		(do-async (p1 "hello"))
		(p1)
	`, s("hello"))

	testCode(t, `
		(def p2 (promise))
		(promise? p1 p2)
	`, a.True)
}
