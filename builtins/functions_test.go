package builtins_test

import (
	"fmt"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {
	as := assert.New(t)

	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("say-hello")
	ns.Delete("identity")

	testCode(t, `
		(defn say-hello
		  "this is a doc string"
		  []
		  "Hello, World!")
		(say-hello)
	`, "Hello, World!")

	testCode(t, `
		(defn identity [value] value)
		(identity "foo")
	`, "foo")

	v, _ := ns.Get("say-hello")
	fv := v.(a.Function)
	as.Equal("this is a doc string", fv.Documentation(), "documented")
}

func TestBadFunction(t *testing.T) {
	symErr := fmt.Sprintf(a.ExpectedSymbol, "99")
	seqErr := fmt.Sprintf(a.ExpectedSequence, "99")
	testBadCode(t, `(defn blah [name 99 bad] (name))`, symErr)
	testBadCode(t, `(defn blah 99 (name))`, seqErr)
	testBadCode(t, `(defn 99 [x y] (+ x y))`, symErr)
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
	testBadCode(t, `(fn 99 "hello")`, fmt.Sprintf(a.ExpectedSequence, "99"))
}

func TestApply(t *testing.T) {
	testCode(t, `(apply + [1 2 3])`, a.NewFloat(6))
	testCode(t, `
		(apply
			(fn [x y z] (+ x y z))
			[1 2 3])
	`, a.NewFloat(6))

	appErr := fmt.Sprintf(a.ExpectedApplicable, "32")
	testBadCode(t, `(apply 32 [1 2 3])`, appErr)
}
