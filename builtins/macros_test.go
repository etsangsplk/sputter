package builtins_test

import "testing"

func TestMacroReplace(t *testing.T) {
	testCode(t, `
        (defmacro foo [] "hello")
        (foo)
    `, "hello")
}
