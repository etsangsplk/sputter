package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestAsync(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("g")

	testCode(t, `
		(def g (async
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
