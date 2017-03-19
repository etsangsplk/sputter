package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestSequence(t *testing.T) {
	testCode(t, `(len '(1, 2, 3))`, a.NewFloat(3))
	testCode(t, `(seq? (list 1 2 3))`, a.True)
	testCode(t, `(seq? ())`, a.False)
	testCode(t, `(first '(1 2 3 4))`, a.NewFloat(1))
	testCode(t, `(first (rest '(1 2 3 4)))`, a.NewFloat(2))
	testCode(t, `(first (rest (cons 1 (list 2 3))))`, a.NewFloat(2))
}

func TestMapFilter(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("x")
	ns.Delete("y")

	testCode(t, `
		(def x '(1 2 3 4 5 6))
		(def y
			(map
				(fn [x] (* x 2))
				(filter (fn [x] (= x 6)) x)))
		(first y)
	`, a.NewFloat(12))
}
