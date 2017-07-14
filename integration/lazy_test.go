package integration_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

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
		(let [ch (chan) emit (:emit ch) close (:close ch) seq (:seq ch)]
			(make-go
				(for-each [i (make-range 1 5 1), j (make-range 1 10 2)]
					(emit (* i j)))
				(close))
			(reduce (lambda [x y] (+ x y)) seq))
	`, f(250))
}
