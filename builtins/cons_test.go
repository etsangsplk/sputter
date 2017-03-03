package builtins_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
)

func TestCons(t *testing.T) {
	testCode(t, `(list? '(1 2 3))`, s.True)
	testCode(t, `(list? ())`, s.True)
	testCode(t, `(list? [1 2 3])`, s.False)
	testCode(t, `(list? 42)`, s.False)
	testCode(t, `(list? (list 1 2 3))`, s.True)
	testCode(t, `(list)`, s.Nil)
	testCode(t, `(first '(1 2 3 4))`, big.NewFloat(1))
	testCode(t, `(first (rest '(1 2 3 4)))`, big.NewFloat(2))
	testCode(t, `(car (cons 1 2))`, big.NewFloat(1))
	testCode(t, `(cdr (cons 1 2))`, big.NewFloat(2))
}

func TestBadCons(t *testing.T) {
	testBadCode(t, `(car 99)`, s.ExpectedCons)
	testBadCode(t, `(cdr 100)`, s.ExpectedCons)
}
