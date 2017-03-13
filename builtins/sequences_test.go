package builtins_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestSequence(t *testing.T) {
	testCode(t, `(len '(1, 2, 3))`, big.NewFloat(3))
	testCode(t, `(seq? (list 1 2 3))`, a.True)
	testCode(t, `(seq? ())`, a.False)
	testCode(t, `(first '(1 2 3 4))`, big.NewFloat(1))
	testCode(t, `(first (rest '(1 2 3 4)))`, big.NewFloat(2))
	testCode(t, `(first (rest (cons 1 (list 2 3))))`, big.NewFloat(2))
}
