package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestMacroReplace(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("hello")
	a.GetNamespace(a.UserDomain).Delete("foo")

	testCode(t, `
		(def hello "you")
		(defmacro foo
			{:doc "this is the macro foo"}
			[x y]
			(+ x y)
			'hello)
		(foo 1 2)
	`, s("you"))
}
