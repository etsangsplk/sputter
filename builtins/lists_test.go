package builtins_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestList(t *testing.T) {
	testCode(t, `(list? '(1 2 3))`, a.True)
	testCode(t, `(list? ())`, a.True)
	testCode(t, `(list? [1 2 3])`, a.False)
	testCode(t, `(list? 42)`, a.False)
	testCode(t, `(list? (list 1 2 3))`, a.True)
	testCode(t, `(list)`, a.EmptyList)
	testCode(t, `(first '(1 2 3 4))`, big.NewFloat(1))
	testCode(t, `(first (rest '(1 2 3 4)))`, big.NewFloat(2))
	testCode(t, `(first (rest (cons 1 (list 2 3))))`, big.NewFloat(2))
}
