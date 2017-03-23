package builtins_test

import (
	"fmt"
	"testing"

	a "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
)

func TestCond(t *testing.T) {
	testCode(t, `(cond)`, a.Nil)

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope"
			true  "hello"
			"hi"  "ignored")
	`, "hello")

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope"
			:else "hello"
			"hi"  "ignored")
	`, "hello")

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope")
	`, a.Nil)
}

func TestBadCond(t *testing.T) {
	error99 := fmt.Sprintf(b.ExpectedCondResult, "99")

	testBadCode(t, `(cond 99)`, error99)

	testBadCode(t, `
		(cond
			true "hello"
			99)
	`, error99)

	testBadCode(t, `
		(cond
			false "hello"
			99)
	`, error99)
}

func TestIf(t *testing.T) {
	testCode(t, `(if false 1 0)`, a.NewFloat(0))
	testCode(t, `(if true 1 0)`, a.NewFloat(1))
	testCode(t, `(if nil 1 0)`, a.NewFloat(0))
	testCode(t, `(if () 1 0)`, a.NewFloat(1))
	testCode(t, `(if "hello" 1 0)`, a.NewFloat(1))
}

func TestBadIfArity(t *testing.T) {
	testBadCode(t, `(if)`, fmt.Sprintf(a.BadArityRange, 2, 3, 0))
}
