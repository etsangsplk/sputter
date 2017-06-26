package core_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	_ "github.com/kode4food/sputter/core"
)

func TestToAssoc(t *testing.T) {
	testCode(t, `(assoc? (to-assoc [:name "Sputter" :age 45]))`, a.True)
	testCode(t, `(assoc? (to-assoc '(:name "Sputter" :age 45)))`, a.True)
	testCode(t, `(mapped? (to-assoc '(:name "Sputter" :age 45)))`, a.True)
}

func TestToVector(t *testing.T) {
	testCode(t, `(vector? (to-vector (list 1 2 3)))`, a.True)
}

func TestToList(t *testing.T) {
	testCode(t, `(list? (to-list (vector 1 2 3)))`, a.True)
}

func TestMapFilter(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("x")
	ns.Delete("y")

	testCode(t, `
		(first (apply list (map (fn [x] (* x 2)) [1 2 3 4])))
	`, f(2))

	testCode(t, `
		(def x (concat '(1 2) (list 3 4)))
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
