package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestBasicNumber(t *testing.T) {
	testCode(t, `(+)`, a.NewFloat(0))
	testCode(t, `(*)`, a.NewFloat(1))
	testCode(t, `(+ 1 1)`, a.NewFloat(2.0))
	testCode(t, `(* 4 4)`, a.NewFloat(16.0))
	testCode(t, `(+ 5 4)`, a.NewFloat(9.0))
	testCode(t, `(* 12 3)`, a.NewFloat(36.0))
	testCode(t, `(- 10 4)`, a.NewFloat(6.0))
	testCode(t, `(- 10 4 2)`, a.NewFloat(4.0))
	testCode(t, `(/ 10 2)`, a.NewFloat(5.0))
	testCode(t, `(/ 10 2 5)`, a.NewFloat(1.0))
}

func TestNestedNumber(t *testing.T) {
	testCode(t, `(/ 10 (- 5 3))`, a.NewFloat(5.0))
	testCode(t, `(* 5 (- 5 3))`, a.NewFloat(10.0))
	testCode(t, `(/ 10 (/ 6 3))`, a.NewFloat(5.0))
}

func TestNonNumber(t *testing.T) {
	helloErr := a.Err(a.ExpectedNumber, `"hello"`)
	testBadCode(t, `(+ 99 "hello")`, helloErr)
	testBadCode(t, `(+ "hello")`, helloErr)
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
	helloErr := a.Err(a.ExpectedNumber, `"hello"`)
	testBadCode(t, `(< 99 "hello")`, helloErr)
	testBadCode(t, `(< "hello" "there")`, helloErr)
}
