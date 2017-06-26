package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
)

func TestDefinitions(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("foo")
	ns.Delete("return-local")

	testCode(t, `
		(def foo "bar")
		foo
	`, s("bar"))

	testCode(t, `
		(def return-local (lambda []
			(let [foo "local"] foo)))
		(return-local)
	`, s("local"))
}

func TestLetBindings(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("foo")
	testCode(t, `
		(def foo 99)
		(let [foo 100]
			(+ foo user:foo))
	`, f(199))

	testBadCode(t, `
		(let 99 "hello")
	`, a.Err(a.ExpectedVector, "99"))

	testBadCode(t, `
		(let [a blah b] "hello")
	`, b.ExpectedBindings)
}
