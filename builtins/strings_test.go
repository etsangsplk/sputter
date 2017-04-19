package builtins_test

import "testing"

func TestStr(t *testing.T) {
	testCode(t, `
	  (str "hello" [1 2 3 4])
	`, "hello[1 2 3 4]")
}
