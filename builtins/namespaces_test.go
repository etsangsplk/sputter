package builtins_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestNamespaces(t *testing.T) {
	testCode(t, `
		(ns foo)
		(def v1 99)
		(ns bar)
		(def v1 100)
		(+ v1 foo:v1)
	`, big.NewFloat(199))

	testBadCode(t, `
		(ns foo:bar)
	`, a.ExpectedUnqualifiedSymbol)
}
