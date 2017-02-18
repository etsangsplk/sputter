package builtins_test

import (
	"math/big"
	"testing"
)

func TestConditional(t *testing.T) {
	testCode(t, `(if false 1 0)`, big.NewFloat(0))
	testCode(t, `(if true 1 0)`, big.NewFloat(1))
	testCode(t, `(if nil 1 0)`, big.NewFloat(0))
	testCode(t, `(if () 1 0)`, big.NewFloat(0))
	testCode(t, `(if "hello" 1 0)`, big.NewFloat(1))
}
