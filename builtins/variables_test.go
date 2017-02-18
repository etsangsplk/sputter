package builtins_test

import "testing"

func TestVariables(t *testing.T) {
	testCode(t, `
		(defvar foo "bar")
		foo
	`, "bar")

	testCode(t, `
		(defun return-local []
			(let [foo "local"] foo))
		(return-local)
	`, "local")
}
