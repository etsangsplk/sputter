package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	"github.com/kode4food/sputter/builtins"
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

	r, _ := ns.Get("say-hello")
	fv := r.(a.Function)
	as.String("this is a doc string", fv.Documentation())
}

func TestBadFunction(t *testing.T) {
	symErr := a.ErrStr(a.ExpectedSymbol, "99")
	vecErr := a.ErrStr(a.ExpectedVector, "99")
	listErr := a.ErrStr(a.ExpectedList, "99")

	testBadCode(t, `(defn blah [name 99 bad] (name))`, symErr)
	testBadCode(t, `(defn blah 99 (name))`, listErr)
	testBadCode(t, `(defn 99 [x y] (+ x y))`, symErr)
	testBadCode(t, `(defn blah (99 "hello"))`, vecErr)
	testBadCode(t, `(defn blah ([x] "hello") 99)`, listErr)
}

func TestBadFunctionArity(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("identity")

	testBadCode(t, `(defn blah)`, a.ErrStr(a.BadMinimumArity, 2, 1))

	testBadCode(t, `
		(defn identity [value] value)
		(identity)
	`, a.ErrStr(builtins.ExpectedArguments, args(local("value"))))
}
