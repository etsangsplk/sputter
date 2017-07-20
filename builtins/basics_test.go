package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	e "github.com/kode4food/sputter/evaluator"
)

func TestDo(t *testing.T) {
	as := assert.New(t)
	c := e.NewEvalContext()

	do := getBuiltIn("do")
	r1 := do(c, args(f(1), f(2), f(3)))
	as.Number(3, r1)
}

func TestReadEval(t *testing.T) {
	as := assert.New(t)
	c := e.NewEvalContext()

	read := getBuiltIn("read")
	eval := getBuiltIn("eval")

	r1 := read(c, args(s("[1 2 3]")))
	e1 := eval(c, args(r1))
	v1 := a.AssertVector(e1)

	v2, ok := v1.ElementAt(0)
	as.True(ok)
	as.Number(1, v2)

	v3, ok := v1.ElementAt(2)
	as.True(ok)
	as.Number(3, v3)
}
