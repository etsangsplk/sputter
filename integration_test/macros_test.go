package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestMacroPredicates(t *testing.T) {
	testCode(t, `(macro? cond)`, a.True)
	testCode(t, `(!macro? cond)`, a.False)
	testCode(t, `(macro? if)`, a.False)
	testCode(t, `(!macro? if)`, a.True)
}

func TestMacroReplace(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("foo")

	testCode(t, `
		(defmacro foo
			{:doc "this is the macro foo"}
			[& args]
			(cons 'str (cons "hello" args)))

		(foo 1 2 3)
	`, s(`hello123`))
}

func TestMacroExpand(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("foo1")
	ns.Delete("foo2")

	testCode(t, `
		(defmacro foo1
			{:doc "this is the macro foo1"}
			[& args]
			(cons 'str (cons "hello" args)))

		(macroexpand1 (foo1 1 2 3))
	`, s(`(str "hello" 1 2 3)`))

	testCode(t, `
		(defmacro foo2
			{:doc "this is the macro foo2"}
			[& args]
			(foo1 (args 0) (args 1) (args 2)))

		(macroexpand (foo2 1 2 3))
	`, s("hello123"))

	testCode(t, `
		(macroexpand-all (foo2 (foo1 1 2 3) 4 5))
	`, s("hello(foo1 1 2 3)45"))
}
