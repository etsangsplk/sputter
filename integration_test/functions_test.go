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
		(def call (lambda [func] (func)))
		(let [greeting "hello"]
			(let [foo (lambda [] greeting)]
				(call foo)))
	`, s("hello"))
}

func TestBadLambda(t *testing.T) {
	err := a.ErrStr(a.ExpectedList, "99")
	testBadCode(t, `(lambda 99 "hello")`, err)

	err = a.ErrStr(a.ExpectedUnqualified, "foo:bar")
	testBadCode(t, `(lambda foo:bar [] "hello")`, err)
}

func TestApply(t *testing.T) {
	testCode(t, `(apply + [1 2 3])`, f(6))
	testCode(t, `
		(apply
			(lambda add {:test true} [x y z] (+ x y z))
			[1 2 3])
	`, f(6))

	appErr := a.ErrStr(a.ExpectedApplicable, "32")
	testBadCode(t, `(apply 32 [1 2 3])`, appErr)
}

func TestRestFunctions(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("test")

	testCode(t, `
		(def test (lambda [f & r] (apply vector (cons f r))))
		(test 1 2 3 4 5 6 7)
	`, a.Str("[1 2 3 4 5 6 7]"))

	testBadCode(t, `
		(lambda [x y &] "explode")
	`, a.ErrStr(builtins.InvalidRestArgument, "[]"))

	testBadCode(t, `
		(lambda [x y & z g] "explode")
	`, a.ErrStr(builtins.InvalidRestArgument, "[z g]"))

	testBadCode(t, `
		(lambda [x y & & z] "explode")
	`, a.ErrStr(builtins.InvalidRestArgument, "[& z]"))
}
