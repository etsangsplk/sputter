package builtins_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
)

func TestCons(t *testing.T) {
	testCode(t, `(list? '(1 2 3))`, s.True)
	testCode(t, `(list? ())`, s.True)
	testCode(t, `(list? [1 2 3])`, s.False)
	testCode(t, `(list? 42)`, s.False)
	testCode(t, `(list? (list 1 2 3))`, s.True)
	testCode(t, `(list)`, s.Nil)
	testCode(t, `(first '(1 2 3 4))`, big.NewFloat(1))
	testCode(t, `(second '(1 2 3 4))`, big.NewFloat(2))
	testCode(t, `(third '(1 2 3 4))`, big.NewFloat(3))
	testCode(t, `(second (cons 1 (cons 2 3)))`, big.NewFloat(2))
	testCode(t, `(car (cons 1 2))`, big.NewFloat(1))
	testCode(t, `(cdr (cons 1 2))`, big.NewFloat(2))
}

func TestBadCons(t *testing.T) {
	testBadCode(t, `(car 99)`, b.NonCons)
	testBadCode(t, `(cdr 100)`, b.NonCons)
}
