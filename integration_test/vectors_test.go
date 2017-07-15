package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestVector(t *testing.T) {
	as := assert.New(t)

	r1 := runCode(`(vector 1 (- 5 3) (+ 1 2))`)
	as.String("[1 2 3]", r1)

	r2 := runCode(`(apply vector (concat '(1) '((- 5 3)) '((+ 1 2))))`)
	as.String("[1 (- 5 3) (+ 1 2)]", r2)

	testCode(t, `(vector? [1 2 3])`, a.True)
	testCode(t, `(vector? (vector 1 2 3))`, a.True)
	testCode(t, `(vector? [])`, a.True)
	testCode(t, `(vector? 99)`, a.False)

	testCode(t, `(!vector? [1 2 3])`, a.False)
	testCode(t, `(!vector? (vector 1 2 3))`, a.False)
	testCode(t, `(!vector? [])`, a.False)
	testCode(t, `(!vector? 99)`, a.True)

	testCode(t, `(len? [1 2 3])`, a.True)
	testCode(t, `(len? 99)`, a.False)
	testCode(t, `(indexed? [1 2 3])`, a.True)
	testCode(t, `(indexed? 99)`, a.False)

	a.GetNamespace(a.UserDomain).Delete("x")
	testCode(t, `
		(def x [1 2 3 4])
		(x 2)
	`, f(3))
}
