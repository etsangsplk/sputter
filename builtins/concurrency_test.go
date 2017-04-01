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

func TestGoConcurrency(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("g")

	testCode(t, `
		(def g (go
			(emit 99)
			(emit 100 1000)))
		(apply + g)
	`, a.NewFloat(1199))
}
