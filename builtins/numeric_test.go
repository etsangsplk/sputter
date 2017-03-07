package builtins_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestBasicNumeric(t *testing.T) {
	testCode(t, `(+ 1 1)`, big.NewFloat(2.0))
	testCode(t, `(* 4 4)`, big.NewFloat(16.0))
	testCode(t, `(+ 5 4)`, big.NewFloat(9.0))
	testCode(t, `(* 12 3)`, big.NewFloat(36.0))
	testCode(t, `(- 10 4)`, big.NewFloat(6.0))
	testCode(t, `(- 10 4 2)`, big.NewFloat(4.0))
	testCode(t, `(/ 10 2)`, big.NewFloat(5.0))
	testCode(t, `(/ 10 2 5)`, big.NewFloat(1.0))
}

func TestNestedNumeric(t *testing.T) {
	testCode(t, `(/ 10 (- 5 3))`, big.NewFloat(5.0))
	testCode(t, `(* 5 (- 5 3))`, big.NewFloat(10.0))
	testCode(t, `(/ 10 (/ 6 3))`, big.NewFloat(5.0))
}

func TestNonNumeric(t *testing.T) {
	testBadCode(t, `(+ 99 "hello")`, a.ExpectedNumeric)
	testBadCode(t, `(+ "hello")`, a.ExpectedNumeric)
}

func TestCompare(t *testing.T) {
	testCode(t, `(= 1 1)`, a.True)
	testCode(t, `(= 1 1 1 1 '1 1 1)`, a.True)
	testCode(t, `(= 1 2)`, a.False)
	testCode(t, `(= 1 1 1 1 2 1 1 1)`, a.False)

	testCode(t, `(!= 1 1)`, a.False)
	testCode(t, `(!= 1 1 1 1 '1 1 1)`, a.False)
	testCode(t, `(!= 1 2)`, a.True)
	testCode(t, `(!= 1 1 1 1 2 1 1 1)`, a.True)

	testCode(t, `(> 1 1)`, a.False)
	testCode(t, `(> 2 1)`, a.True)
	testCode(t, `(> 1 2)`, a.False)
	testCode(t, `(> 1 2 3 4 5)`, a.False)
	testCode(t, `(> 5 4 3 2 1)`, a.True)
	testCode(t, `(>= 1 1)`, a.True)
	testCode(t, `(>= 0 1)`, a.False)
	testCode(t, `(>= 1 0)`, a.True)

	testCode(t, `(< 1 1)`, a.False)
	testCode(t, `(< 2 1)`, a.False)
	testCode(t, `(< 1 2)`, a.True)
	testCode(t, `(< 1 2 3 4 5)`, a.True)
	testCode(t, `(< 5 4 3 2 1)`, a.False)
	testCode(t, `(<= 1 1)`, a.True)
	testCode(t, `(<= 0 1)`, a.True)
	testCode(t, `(<= 1 0)`, a.False)
}

func TestBadCompare(t *testing.T) {
	testBadCode(t, `(< 99 "hello")`, a.ExpectedNumeric)
	testBadCode(t, `(< "hello" "there")`, a.ExpectedNumeric)
}
