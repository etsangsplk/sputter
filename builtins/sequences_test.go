package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestSequence(t *testing.T) {
	testCode(t, `(len '(1, 2, 3))`, f(3))
	testCode(t, `(len? '(1 2 3))`, a.True)
	testCode(t, `(len? (make-range 1 5 1))`, a.False)
	testCode(t, `(indexed? '(1 2 3))`, a.True)
	testCode(t, `(indexed? [1 2 3])`, a.True)
	testCode(t, `(indexed? (make-range 1 5 1))`, a.False)

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

func TestRange(t *testing.T) {
	testCode(t, `
		(reduce
			(lambda [x y] (+ x y))
			(make-range 1 5 1))
	`, f(10))

	testCode(t, `
		(reduce
			(lambda [x y] (+ x y))
			(make-range 5 1 -1))
	`, f(14))
}

func TestMapAndFilter(t *testing.T) {
	testCode(t, `
		(reduce
			(lambda [x y] (+ x y))
			(map
				(lambda [x] (* x 2))
				(filter
					(lambda [x] (<= x 5))
					[1 2 3 4 5 6 7 8 9 10])))
	`, f(30))
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
		(nth (apply vector (take 6 x y)) 5)
	`, f(6))

	ns.Delete("x")
	ns.Delete("y")
	testCode(t, `
		(def x '(1 2 3 4))
		(def y [5 6 7 8])
		(nth (apply vector (drop 3 x y)) 0)
	`, f(4))
}

func TestForEachLoop(t *testing.T) {
	testCode(t, `
		(let [ch (channel) emit (:emit ch) close (:close ch) seq (:seq ch)]
			(do-async
				(for-each [i (make-range 1 5 1), j (make-range 1 10 2)]
					(emit (* i j)))
				(close))
			(reduce (lambda [x y] (+ x y)) seq))
	`, f(250))
}
