package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestHashMap(t *testing.T) {
	testCode(t, `(len {:name "Sputter", :age 45})`, a.NewFloat(2))
	testCode(t, `(len (hash-map :name "Sputter", :age 45))`, a.NewFloat(2))
	testCode(t, `(hash-map? {:name "Sputter" :age 45})`, a.True)
	testCode(t, `(hash-map? (hash-map :name "Sputter" :age 45))`, a.True)
	testCode(t, `(hash-map? '(:name "Sputter" :age 45))`, a.False)
	testCode(t, `(hash-map? (to-hash-map '(:name "Sputter" :age 45)))`, a.True)
	testCode(t, `(hash-map? [:name "Sputter" :age 45])`, a.False)
	testCode(t, `(!hash-map? '(:name "Sputter" :age 45))`, a.True)
	testCode(t, `(!hash-map? [:name "Sputter" :age 45])`, a.True)
	testCode(t, `(:name {:name "Sputter" :age 45})`, "Sputter")

	testCode(t, `
		(:name (apply hash-map (concat '(:name "Sputter") '(:age 45))))
	`, "Sputter")

	a.GetNamespace(a.UserDomain).Delete("x")
	testCode(t, `
		(def x {:name "bob" :age 45})
		(x :name)
	`, "bob")
	
	testBadCode(t, `(hash-map :too "few" :args)`, a.ExpectedPair)

	testBadCode(t, `
		(apply hash-map (concat '(:name "Sputter") '(:age)))
	`, a.ExpectedPair)
}
