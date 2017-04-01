package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestGoConcurrency(t *testing.T) {
	testCode(t, `
		(def x (go
			(emit 99)
			(emit 100 1000)))
		(def r (to-vector x))
		(+ (first r) (first (rest r)) (first (rest (rest r))))
	`, a.NewFloat(1199))
}
