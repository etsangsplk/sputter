package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestVector(t *testing.T) {
	as := assert.New(t)

	r1 := runCode(`(vector 1 (- 5 3) (+ 1 2))`)
	as.Equal("[1 2 3]", a.String(r1), "correct vector")

	r2 := runCode(`(apply vector (concat '(1) '((- 5 3)) '((+ 1 2))))`)
	as.Equal("[1 2 3]", a.String(r2), "correct vector")

	testCode(t, `(vector? [1 2 3])`, a.True)
	testCode(t, `(vector? (vector 1 2 3))`, a.True)
	testCode(t, `(vector? [])`, a.True)
	testCode(t, `(vector? 99)`, a.False)

	testCode(t, `(!vector? [1 2 3])`, a.False)
	testCode(t, `(!vector? (vector 1 2 3))`, a.False)
	testCode(t, `(!vector? [])`, a.False)
	testCode(t, `(!vector? 99)`, a.True)

	testCode(t, `(vector? (to-vector (list 1 2 3)))`, a.True)

	a.GetNamespace(a.UserDomain).Delete("x")
	testCode(t, `
		(def x [1 2 3 4])
		(x 2)
	`, a.NewFloat(3))	
}
