package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	_ "github.com/kode4food/sputter/core"
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

func runCodeWithContext(c a.Context, src string) a.Value {
	return e.EvalStr(c, a.Str(src))
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

func TestBuiltInsContext(t *testing.T) {
	as := assert.New(t)

	bg1 := e.NewEvalContext()
	bg2 := a.ChildContext(bg1)
	bg3 := a.ChildContext(bg2)

	qv, ok := bg3.Get("do")
	as.True(ok)
	if _, ok := qv.(a.Applicable); !ok {
		as.Fail("returned value not Applicable")
	}
}

func TestDo(t *testing.T) {
	testCode(t, `
		(do
			55
			(if true 99 33))
	`, f(99))
}

func TestTrueFalse(t *testing.T) {
	testCode(t, `true`, a.True)
	testCode(t, `false`, a.False)
	testCode(t, `nil`, a.Nil)
}

func TestReadEval(t *testing.T) {
	testCode(t, `
		(eval (read "(str \"hello\" \"you\" \"test\")"))
	`, s("helloyoutest"))
}
