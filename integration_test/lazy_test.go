package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestRange(t *testing.T) {
	testCode(t, `
		(reduce
			(fn [x y] (+ x y))
			(range* 1 5 1))
	`, f(10))

	testCode(t, `
		(reduce
			(fn [x y] (+ x y))
			(range* 5 1 -1))
	`, f(14))
}

func TestMapAndFilter(t *testing.T) {
	testCode(t, `
		(reduce
			(fn [x y] (+ x y))
			(map
				(fn [x] (* x 2))
				(filter
					(fn [x] (<= x 5))
					[1 2 3 4 5 6 7 8 9 10])))
	`, f(30))
}

func TestMapParallel(t *testing.T) {
	testCode(t, `
		(to-vector
			(map + 
				[1 2 3 4]
				'(2 4 6 8)
				(range 20 30)))
	`, s("[23 27 31 35]"))
}

func TestReduce(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("x")
	ns.Delete("y")

	testCode(t, `
		(def x '(1 2 3 4))
		(reduce + x)
	`, f(10))

	testCode(t, `
		(def y (concat '(1 2 3 4) [5 6 7 8]))
		(reduce + y)
	`, f(36))

	testCode(t, `
		(reduce + 10 y)
	`, f(46))
}

func TestTakeDrop(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)

	ns.Delete("x")
	testCode(t, `
		(def x (concat '(1 2 3 4) [5 6 7 8]))
		(nth (apply vector (take 6 x)) 5)
	`, f(6))

	ns.Delete("x")
	testCode(t, `
		(def x (concat '(1 2 3 4) [5 6 7 8]))
		(nth (apply vector (drop 3 x)) 0)
	`, f(4))
}

func TestLazySeq(t *testing.T) {
	testCode(t, `
		(reduce
			(fn [x y] (+ x y))
			(lazy-seq (cons 1 (lazy-seq [2, 3]))))
	`, f(6))

	testCode(t, `
		(len (to-vector (lazy-seq nil)))
	`, f(0))
}

func TestForEachLoop(t *testing.T) {
	testCode(t, `
		(let [ch (chan) emit (:emit ch) close (:close ch) seq (:seq ch)]
			(go*
				(for-each [i (range* 1 5 1), j (range* 1 10 2)]
					(emit (* i j)))
				(close))
			(reduce (fn [x y] (+ x y)) seq))
	`, f(250))
}
