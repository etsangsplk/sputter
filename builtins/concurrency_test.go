package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestChannel(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("c")
	ns.Delete("e")
	ns.Delete("s")
	
	testCode(t, `
		(def c (channel))
		(def e (:emit c))
		(def s (:seq c))
		
		(go (e "hello"))	
		(first s)
	`, "hello")
}

func TestGoConcurrency(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("g")
	ns.Delete("r")

	testCode(t, `
		(def g (go
			(emit 99)
			(emit 100 1000)))
		(def r (to-vector g))
		(+ (first r) (first (rest r)) (first (rest (rest r))))
	`, a.NewFloat(1199))
}
