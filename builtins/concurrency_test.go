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
		(async
			(apply (:emit c) '("hello"))
			(apply (:close c) ()))
		(first (:seq c))
	`, s("hello"))
}

func TestLazySequence(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("g")

	testCode(t, `
		(def g (generate
			(emit 99)
			(emit 100 1000)))
		(apply + g)
	`, f(1199))
}

func TestPromise(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("p")

	testCode(t, `
		(def p (promise))
		(async (p "hello"))
		(p)
	`, s("hello"))
}

func TestFuture(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("p")

	testCode(t, `
		(def p (future "hello"))
		(p)
	`, s("hello"))
}
