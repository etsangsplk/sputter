package integration_test

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

func TestFuture(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("p")

	testCode(t, `
		(def p (future "hello"))
		(p)
	`, s("hello"))
}
