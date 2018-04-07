package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/builtins"
)

func TestFunctionPredicates(t *testing.T) {
	testCode(t, `(apply? if)`, a.True)
	testCode(t, `(!apply? if)`, a.False)
	testCode(t, `(apply? 99)`, a.False)
	testCode(t, `(!apply? 99)`, a.True)

	testCode(t, `(special-form? if)`, a.True)
	testCode(t, `(!special-form? if)`, a.False)
	testCode(t, `(special-form? eq)`, a.False)
	testCode(t, `(!special-form? eq)`, a.True)
}

func TestLambda(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("call")

	testCode(t, `
		(def call (fn [func] (func)))
		(let [greeting "hello"]
			(let [foo (fn [] greeting)]
				(call foo)))
	`, s("hello"))
}

func TestBadLambda(t *testing.T) {
	e := typeErr("*api.dec", "*api.List")
	testBadCode(t, `(fn 99 "hello")`, e)

	e = intfErr("*api.qualifiedSymbol", "api.LocalSymbol", "LocalSymbolType")
	testBadCode(t, `(fn foo:bar [] "hello")`, e)
}

func TestApply(t *testing.T) {
	testCode(t, `(apply + [1 2 3])`, f(6))
	testCode(t, `
		(apply
			(fn add {:test true} [x y z] (+ x y z))
			[1 2 3])
	`, f(6))

	e := intfErr("*api.dec", "api.Applicable", "Apply")
	testBadCode(t, `(apply 32 [1 2 3])`, e)
}

func TestRestFunctions(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("test")

	testCode(t, `
		(def test (fn [f & r] (apply vector (cons f r))))
		(test 1 2 3 4 5 6 7)
	`, a.Str("[1 2 3 4 5 6 7]"))

	testBadCode(t, `
		(fn [x y &] "explode")
	`, a.ErrStr(builtins.InvalidRestArgument, "[]"))

	testBadCode(t, `
		(fn [x y & z g] "explode")
	`, a.ErrStr(builtins.InvalidRestArgument, "[z g]"))

	testBadCode(t, `
		(fn [x y & & z] "explode")
	`, a.ErrStr(builtins.InvalidRestArgument, "[& z]"))
}
