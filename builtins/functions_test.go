package builtins_test

import (
	"fmt"
	"testing"

	b "github.com/kode4food/sputter/builtins"
)

func TestFunction(t *testing.T) {
	testCode(t, `
		(defun say-hello [] "Hello, World!")
		(say-hello)
	`, "Hello, World!")

	testCode(t, `
		(defun identity [value] value)
		(identity "foo")
	`, "foo")
}

func TestBadFunctionArity(t *testing.T) {
	testBadCode(t, `(defun blah)`, fmt.Sprintf(b.BadMinArity, 3, 1))

	testBadCode(t, `
		(defun identity [value] value)
		(identity)
	`, fmt.Sprintf(b.BadArity, 1, 0))
}
