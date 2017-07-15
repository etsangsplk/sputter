package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestGenerate(t *testing.T) {
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
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("p1")
	ns.Delete("p2")

	testCode(t, `
		(def p1 (promise))
		(promise? p1)
	`, a.True)

	testCode(t, `
		(def p2 (promise "hello"))
		(p2)
	`, s("hello"))
}

func TestFuture(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("p")

	testCode(t, `
		(def p (future "hello"))
		(p)
	`, s("hello"))
}
