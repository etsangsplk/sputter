package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestMap(t *testing.T) {
	testCode(t, `(len {:name "Sputter", :age 45})`, a.NewFloat(2))
	testCode(t, `(len (map :name "Sputter", :age 45))`, a.NewFloat(2))
	testCode(t, `(map? {:name "Sputter" :age 45})`, a.True)
	testCode(t, `(map? (map :name "Sputter" :age 45))`, a.True)
	testCode(t, `(map? '(:name "Sputter" :age 45))`, a.False)
	testCode(t, `(map? (to-map '(:name "Sputter" :age 45)))`, a.True)
	testCode(t, `(map? [:name "Sputter" :age 45])`, a.False)
	testCode(t, `(!map? '(:name "Sputter" :age 45))`, a.True)
	testCode(t, `(!map? [:name "Sputter" :age 45])`, a.True)
	testCode(t, `(:name {:name "Sputter" :age 45})`, "Sputter")

	testBadCode(t, `(map :too "few" :args)`, a.ExpectedPair)
}
