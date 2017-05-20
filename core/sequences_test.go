package core_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
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
