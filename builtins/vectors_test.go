package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestNewVector(t *testing.T) {
	as := assert.New(t)
	r := runCode(`(vector 1 (- 5 3) (+ 1 2))`)
	as.Equal("[1 2 3]", a.String(r), "correct vector")

	testCode(t, `(vector? [1 2 3])`, a.True)
	testCode(t, `(vector? (vector 1 2 3))`, a.True)
	testCode(t, `(vector? [])`, a.True)
	testCode(t, `(vector? 99)`, a.False)

	testCode(t, `(!vector? [1 2 3])`, a.False)
	testCode(t, `(!vector? (vector 1 2 3))`, a.False)
	testCode(t, `(!vector? [])`, a.False)
	testCode(t, `(!vector? 99)`, a.True)
}
