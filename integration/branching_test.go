package integration_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestCond(t *testing.T) {
	testCode(t, `(cond)`, a.Nil)

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope"
			true  "hello"
			"hi"  "ignored")
	`, s("hello"))

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope"
			:else "hello"
			"hi"  "ignored")
	`, s("hello"))

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope")
	`, a.Nil)

	testCode(t, `
		(cond
			true "hello"
			99)
	`, s("hello"))

	testCode(t, `(cond 99)`, f(99))

	testCode(t, `
		(cond
			false "hello"
			99)
	`, f(99))
}
