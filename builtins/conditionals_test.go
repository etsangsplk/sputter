package builtins_test

import (
	"fmt"
	"math/big"
	"testing"

	b "github.com/kode4food/sputter/builtins"
)

func TestConditional(t *testing.T) {
	testCode(t, `(if false 1 0)`, big.NewFloat(0))
	testCode(t, `(if true 1 0)`, big.NewFloat(1))
	testCode(t, `(if nil 1 0)`, big.NewFloat(0))
	testCode(t, `(if () 1 0)`, big.NewFloat(0))
	testCode(t, `(if "hello" 1 0)`, big.NewFloat(1))
}

func TestBadIfArity(t *testing.T) {
	testBadCode(t, `(if)`, fmt.Sprintf(b.BadArityRange, 2, 3, 0))
}
