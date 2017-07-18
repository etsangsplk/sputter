package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	b "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

func getBuiltIn(n a.Name) a.SequenceProcessor {
	if r, ok := b.GetBuiltIn(n); ok {
		return r
	}
	panic(a.Err("Built in not found: ", n))
}

func f(n float64) a.Number {
	return a.NewFloat(n)
}

func args(e ...a.Value) a.Vector {
	return a.NewVector(e...)
}

func s(v string) a.Str {
	return a.Str(v)
}

func kw(n a.Name) a.Keyword {
	return a.NewKeyword(n)
}

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

func TestMissingBuiltIn(t *testing.T) {
	as := assert.New(t)

	v, ok := b.GetBuiltIn("boom")
	as.Nil(v)
	as.False(ok)
}

func TestExplodingBuiltInCall(t *testing.T) {
	as := assert.New(t)
	read := getBuiltIn("read")
	eval := getBuiltIn("eval")

	defer as.ExpectError(a.Err(a.KeyNotFound, a.Name("boom")))
	c := e.NewEvalContext()
	r := read(c, args(s("(def-builtin boom)")))
	eval(c, args(r))
}
