package builtins_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestMap(t *testing.T) {
	testCode(t, `(len {:name "Sputter", :age 45})`, big.NewFloat(2))
	testCode(t, `(len (map :name "Sputter", :age 45))`, big.NewFloat(2))

	testBadCode(t, `(map :too "few" :args)`, a.ExpectedPair)
}
