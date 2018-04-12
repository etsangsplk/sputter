package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestBasicNumber(t *testing.T) {
	testCode(t, `(+)`, f(0))
	testCode(t, `(*)`, f(1))
	testCode(t, `(+ 1 1)`, f(2.0))
	testCode(t, `(* 4 4)`, f(16.0))
	testCode(t, `(+ 5 4)`, f(9.0))
	testCode(t, `(* 12 3)`, f(36.0))
	testCode(t, `(- 10 4)`, f(6.0))
	testCode(t, `(- 10 4 2)`, f(4.0))
	testCode(t, `(/ 10 2)`, f(5.0))
	testCode(t, `(/ 10 2 5)`, f(1.0))
	testCode(t, `(mod 10 3)`, f(1.0))
	testCode(t, `(mod 100 8 7)`, f(4.0))
}

func TestNestedNumber(t *testing.T) {
	testCode(t, `(/ 10 (- 5 3))`, f(5.0))
	testCode(t, `(* 5 (- 5 3))`, f(10.0))
	testCode(t, `(/ 10 (/ 6 3))`, f(5.0))
}

func TestNonNumber(t *testing.T) {
	e := intfErr("api.Str", "api.Number", "Add")
	testBadCode(t, `(+ 99 "hello")`, e)
	testBadCode(t, `(+ "hello")`, e)
}

func TestBadNumbers(t *testing.T) {
	testBadNumber := func(err string, ns string) {
		testBadCode(t, ns, a.ErrStr(err, a.Str(ns)))
	}

	testBadNumber(a.ExpectedInteger, "0xfkk")
	testBadNumber(a.ExpectedInteger, "0b01109")
	testBadNumber(a.ExpectedInteger, "123j-k")
	testBadNumber(a.ExpectedFloat, "1.2j-k")
	testBadNumber(a.ExpectedRatio, "1/2p")
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
	e := intfErr("api.Str", "api.Number", "Add")
	testBadCode(t, `(< 99 "hello")`, e)
	testBadCode(t, `(< "hello" "there")`, e)
}
