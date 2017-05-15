package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestFunction(t *testing.T) {
	as := assert.New(t)

	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("say-hello")
	ns.Delete("identity")

	testCode(t, `
		(defn say-hello
		  "this is a doc string"
		  {:test true}
		  []
		  "Hello, World!")
		(say-hello)
	`, s("Hello, World!"))

	testCode(t, `
		(defn identity [value] value)
		(identity "foo")
	`, s("foo"))

	v, _ := ns.Get("say-hello")
	fv := v.(a.Function)
	as.String("this is a doc string", fv.Documentation())
}

func TestBadFunction(t *testing.T) {
	symErr := a.Err(a.ExpectedSymbol, "99")
	vecErr := a.Err(a.ExpectedVector, "99")
	testBadCode(t, `(defn blah [name 99 bad] (name))`, symErr)
	testBadCode(t, `(defn blah 99 (name))`, vecErr)
	testBadCode(t, `(defn 99 [x y] (+ x y))`, symErr)
}

func TestBadFunctionArity(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("identity")

	testBadCode(t, `(defn blah)`, a.Err(a.BadMinimumArity, 3, 1))

	testBadCode(t, `
		(defn identity [value] value)
		(identity)
	`, a.Err(a.BadArity, 1, 0))
}

func TestLambda(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("call")

	testCode(t, `
		(defn call [func] (func))
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
