package core_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	e "github.com/kode4food/sputter/evaluator"
)

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) a.Number {
	return a.NewFloat(f)
}

func runCode(src string) a.Value {
	return e.EvalStr(e.NewEvalContext(), a.Str(src))
}

func testCode(t *testing.T, src string, expect a.Value) {
	as := assert.New(t)
	as.Equal(expect, runCode(src))
}

func testBadCode(t *testing.T, src string, err string) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.String(err, rec)
			return
		}
		as.Fail("bad code should panic")
	}()

	runCode(src)
}

func TestCond(t *testing.T) {
	testCode(t, `(cond)`, a.Nil)

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope"
			true  "hello"
			"hi"  "ignored")
	`, s("hello"))

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope"
			:else "hello"
			"hi"  "ignored")
	`, s("hello"))

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope")
	`, a.Nil)

	testCode(t, `
		(cond
			true "hello"
			99)
	`, s("hello"))

	testCode(t, `(cond 99)`, f(99))

	testCode(t, `
		(cond
			false "hello"
			99)
	`, f(99))
}
