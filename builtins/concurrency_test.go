package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestChannel(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("c")

	testCode(t, `
		(def c (channel 1))
		(apply (:emit c) '("hello")) ; buffer of 1
		(apply (:close c) ())
		(first (:seq c))
	`, "hello")
}

func TestAsync(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("g")

	testCode(t, `
		(def g (async
			(emit 99)
			(emit 100 1000)))
		(apply + g)
	`, a.NewFloat(1199))
}

func TestPromise(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("p")

	testCode(t, `
		(def p (promise))
		(async (p "hello"))
		(p)
	`, "hello")
}

func TestFuture(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("p")

	testCode(t, `
		(def p (future "hello"))
		(p)
	`, "hello")
}
