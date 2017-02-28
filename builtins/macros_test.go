package builtins_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
)

func TestMacroReplace(t *testing.T) {
	s.GetNamespace(s.UserDomain).Delete("foo")
	testCode(t, `
        (defmacro foo [] "hello")
        (foo)
    `, "hello")
}
