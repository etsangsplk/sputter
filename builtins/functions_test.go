package builtins_test

import (
	"fmt"
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestFunction(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("say-hello")
	ns.Delete("identity")

	testCode(t, `
		(defn say-hello [] "Hello, World!")
		(say-hello)
	`, "Hello, World!")

	testCode(t, `
		(defn identity [value] value)
		(identity "foo")
	`, "foo")
}

func TestBadFunction(t *testing.T) {
	testBadCode(t, `(defn blah [name 99 bad] (name))`, a.ExpectedUnqualified)
	testBadCode(t, `(defn blah 99 (name))`, a.ExpectedSequence)
	testBadCode(t, `(defn 99 [x y] (+ x y))`, a.ExpectedUnqualified)
}

func TestBadFunctionArity(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("identity")

	testBadCode(t, `(defn blah)`, fmt.Sprintf(a.BadMinimumArity, 3, 1))

	testBadCode(t, `
		(defn identity [value] value)
		(identity)
	`, fmt.Sprintf(a.BadArity, 1, 0))
}

func TestLambda(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("call")

	testCode(t, `
		(defn call [func] (func))
		(let [greeting "hello"]
			(let [foo (fn [] greeting)]
				(call foo)))
	`, "hello")
}

func TestBadLambda(t *testing.T) {
	testBadCode(t, `(fn 99 "hello")`, a.ExpectedSequence)
}

func TestApply(t *testing.T) {
	testCode(t, `(apply + [1 2 3])`, big.NewFloat(6))
	testCode(t, `
		(apply
			(fn [x y z] (+ x y z))
			[1 2 3])
	`, big.NewFloat(6))

	testBadCode(t, `(apply 32 [1 2 3])`, a.ExpectedApplicable)
}
