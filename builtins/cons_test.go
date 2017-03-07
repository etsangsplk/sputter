package builtins_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestCons(t *testing.T) {
	testCode(t, `(list? '(1 2 3))`, a.True)
	testCode(t, `(list? ())`, a.True)
	testCode(t, `(list? [1 2 3])`, a.False)
	testCode(t, `(list? 42)`, a.False)
	testCode(t, `(list? (list 1 2 3))`, a.True)
	testCode(t, `(list)`, a.Nil)
	testCode(t, `(first '(1 2 3 4))`, big.NewFloat(1))
	testCode(t, `(first (rest '(1 2 3 4)))`, big.NewFloat(2))
	testCode(t, `(car (cons 1 2))`, big.NewFloat(1))
	testCode(t, `(cdr (cons 1 2))`, big.NewFloat(2))
}

func TestBadCons(t *testing.T) {
	testBadCode(t, `(car 99)`, a.ExpectedCons)
	testBadCode(t, `(cdr 100)`, a.ExpectedCons)
}
