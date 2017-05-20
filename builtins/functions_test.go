package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/builtins"
)

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
	err := a.Err(a.ExpectedVector, "99")
	testBadCode(t, `(fn 99 "hello")`, err)

	err = a.Err(a.ExpectedUnqualified, "foo:bar")
	testBadCode(t, `(fn foo:bar [] "hello")`, err)
}

func TestApply(t *testing.T) {
	testCode(t, `(apply + [1 2 3])`, f(6))
	testCode(t, `
		(apply
			(fn add {:test true} [x y z] (+ x y z))
			[1 2 3])
	`, f(6))

	appErr := a.Err(a.ExpectedApplicable, "32")
	testBadCode(t, `(apply 32 [1 2 3])`, appErr)
}

func TestRestFunctions(t *testing.T) {
	testCode(t, `
		(def test (fn [f & r] (apply vector (cons f r))))
		(test 1 2 3 4 5 6 7)
	`, a.Str("[1 2 3 4 5 6 7]"))

	testBadCode(t, `
		(fn [x y &] "explode")
	`, a.Err(builtins.InvalidRestArgument, "[&]"))

	testBadCode(t, `
		(fn [x y & z g] "explode")
	`, a.Err(builtins.InvalidRestArgument, "[& z g]"))

	testBadCode(t, `
		(fn [x y & & z] "explode")
	`, a.Err(builtins.InvalidRestArgument, "[& & z]"))
}
