package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestSequence(t *testing.T) {
	testCode(t, `(len '(1, 2, 3))`, f(3))
	testCode(t, `(seq? (list 1 2 3))`, a.True)
	testCode(t, `(seq? ())`, a.False)
	testCode(t, `(seq? (list 1 2 3) 99)`, a.False)
	testCode(t, `(first '(1 2 3 4))`, f(1))
	testCode(t, `(first (rest '(1 2 3 4)))`, f(2))
	testCode(t, `(first (rest (cons 1 (list 2 3))))`, f(2))
	testCode(t, `(first (rest (conj (list 2 3) 1)))`, f(2))

	testCode(t, `(nth '(1 2 3) 1)`, f(2))
	testCode(t, `(nth '(1 2 3) 5 "nope")`, s("nope"))
	testBadCode(t, `(nth '(1 2 3) 5)`, a.Err(a.IndexNotFound, "5"))
}

func TestMapFilter(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("x")
	ns.Delete("y")

	testCode(t, `
	    ;; the conj below is invalid, but never evaluated
		(def x (concat '(1 2) (conj 3 (list 4))))
		(def y
			(map
				(fn [x] (* x 2))
				(filter
					(fn [x] (= x 6))
					[5 6])))
		(apply +
			(map
				(fn [z] (first z))
				[x y]))
	`, f(13))
}

func TestReduce(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("x")
	ns.Delete("y")

	testCode(t, `
		(def x '(1 2 3 4))
		(def y [5 6 7 8])
		(reduce + x y)
	`, f(36))
}

func TestTakeDrop(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)

	ns.Delete("x")
	ns.Delete("y")
	testCode(t, `
		(def x '(1 2 3 4))
		(def y [5 6 7 8])
		(nth (to-vector (take 6 x y)) 5)
	`, f(6))

	ns.Delete("x")
	ns.Delete("y")
	testCode(t, `
		(def x '(1 2 3 4))
		(def y [5 6 7 8])
		(nth (to-vector (drop 3 x y)) 0)
	`, f(4))
}
