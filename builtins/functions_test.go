package builtins_test

import (
	"fmt"
	"testing"

	s "github.com/kode4food/sputter/api"
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

func TestBadFunction(t *testing.T) {
	testBadCode(t, `(defun blah [name 99 bad] (name))`, s.ExpectedSymbol)
	testBadCode(t, `(defun blah 99 (name))`, s.ExpectedSequence)
	testBadCode(t, `(defun 99 [x y] (+ x y))`, s.ExpectedSymbol)
}

func TestBadFunctionArity(t *testing.T) {
	testBadCode(t, `(defun blah)`, fmt.Sprintf(s.BadMinimumArity, 3, 1))

	testBadCode(t, `
		(defun identity [value] value)
		(identity)
	`, fmt.Sprintf(s.BadArity, 1, 0))
}

func TestLambda(t *testing.T) {
	testCode(t, `
		(defun call [func] (func))
		(let [greeting "hello"]
			(let [foo (lambda [] greeting)]
				(call foo)))
	`, "hello")
}

func TestBadLambda(t *testing.T) {
	testBadCode(t, `(lambda 99 "hello")`, s.ExpectedSequence)
}
