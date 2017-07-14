package integration_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestIf(t *testing.T) {
	testCode(t, `(if false 1 0)`, f(0))
	testCode(t, `(if true 1 0)`, f(1))
	testCode(t, `(if nil 1 0)`, f(0))
	testCode(t, `(if () 1 0)`, f(1))
	testCode(t, `(if "hello" 1 0)`, f(1))
	testCode(t, `(if false 1)`, a.Nil)
}

func TestBadIfArity(t *testing.T) {
	testBadCode(t, `(if)`, a.Err(a.BadArityRange, 2, 3, 0))
}
