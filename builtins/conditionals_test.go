package builtins_test

import (
	"fmt"
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
)

func TestCond(t *testing.T) {
	testCode(t, `(cond)`, s.Nil)

	testCode(t, `
		(cond
			(false "goodbye")
			(nil   "nope")
			(true  "hello")
			("hi"  "ignored"))
	`, "hello")

	testCode(t, `
		(cond
			(false "goodbye")
			(nil   "nope"))
	`, s.Nil)

	testCode(t, `
		(cond
			(true "hello")
			99)
	`, "hello")
}

func TestBadCond(t *testing.T) {
	testBadCode(t, `(cond 99)`, s.ExpectedSequence)

	testBadCode(t, `
		(cond
			(false "hello")
			99)
	`, s.ExpectedSequence)
}

func TestIf(t *testing.T) {
	testCode(t, `(if false 1 0)`, big.NewFloat(0))
	testCode(t, `(if true 1 0)`, big.NewFloat(1))
	testCode(t, `(if nil 1 0)`, big.NewFloat(0))
	testCode(t, `(if () 1 0)`, big.NewFloat(0))
	testCode(t, `(if "hello" 1 0)`, big.NewFloat(1))
}

func TestBadIfArity(t *testing.T) {
	testBadCode(t, `(if)`, fmt.Sprintf(s.BadArityRange, 2, 3, 0))
}
