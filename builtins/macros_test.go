package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestMacroReplace(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("foo")

	runCode(`
		(defmacro foo
			{:doc "this is the macro foo"}
			[aList]
			aList)
	`)

	testCode(t, `
		(foo (1 2 3))
	`, s("(1 2 3)"))
}
