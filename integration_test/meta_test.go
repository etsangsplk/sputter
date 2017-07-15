package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestMeta(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("x")
	testCode(t, `
		(def x (with-meta (lambda [] "hello") {
			:foo "bar"}))
		(:foo (meta x))
	`, s("bar"))

	testCode(t, `(meta? (lambda [x] 1))`, a.True)
	testCode(t, `(meta? if)`, a.True)
	testCode(t, `(meta? 99)`, a.False)
	testCode(t, `(meta? "hello")`, a.False)

	testCode(t, `(!meta? (lambda [x] 1))`, a.False)
	testCode(t, `(!meta? if)`, a.False)
	testCode(t, `(!meta? 99)`, a.True)
	testCode(t, `(!meta? "hello")`, a.True)
}
