package integration_test_test

import (
	"testing"

	"fmt"
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

func v(e ...a.Value) a.Vector {
	return a.NewVector(e...)
}

func args(e ...a.Value) a.Vector {
	return v(e...)
}

func local(n a.Name) a.Symbol {
	return a.NewLocalSymbol(n)
}

func cvtErr(concrete, intf, method string) a.Error {
	err := "interface conversion: %s is not %s: missing method %s"
	return a.ErrStr(fmt.Sprintf(err, concrete, intf, method))
}

func runCode(src string) a.Value {
	return e.EvalStr(e.NewEvalContext(), a.Str(src))
}

func testCode(t *testing.T, src string, expect a.Value) {
	as := assert.New(t)
	as.Equal(expect, runCode(src))
}

func testBadCode(t *testing.T, src string, err error) {
	as := assert.New(t)

	defer as.ExpectError(err)
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
