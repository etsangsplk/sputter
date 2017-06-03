package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

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
	a.GetNamespace(a.UserDomain).Delete("foo")

	testCode(t, `
		(defmacro foo
			{:doc "this is the macro foo"}
			[& args]
			(cons 'str (cons "hello" args)))

		(macroexpand (foo 1 2 3))
	`, s(`(str "hello" 1 2 3)`))
}
