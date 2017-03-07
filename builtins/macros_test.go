package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestMacroReplace(t *testing.T) {
	a.GetNamespace(a.UserDomain).Delete("foo")
	testCode(t, `
        (defmacro foo [] "hello")
        (foo)
    `, "hello")
}
