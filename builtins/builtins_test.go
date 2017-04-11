package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	_ "github.com/kode4food/sputter/builtins"
	p "github.com/kode4food/sputter/parser"

	"github.com/stretchr/testify/assert"
)

func runCode(src string) a.Value {
	l := p.NewLexer(src)
	c := a.NewEvalContext()
	tr := p.NewReader(c, l)
	return p.EvalReader(c, tr)
}

func testCode(t *testing.T, src string, expect a.Value) {
	as := assert.New(t)
	if expectNum, ok := expect.(*a.Number); ok {
		as.Equal(a.EqualTo, expectNum.Cmp(runCode(src).(*a.Number)))
		return
	}
	as.Equal(expect, runCode(src), src)
}

func testBadCode(t *testing.T, src string, err string) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.Equal(err, rec, "bad code panics properly")
			return
		}
		as.Fail("bad code should panic")
	}()

	runCode(src)
}

func TestBuiltInsContext(t *testing.T) {
	as := assert.New(t)

	bg1 := a.NewEvalContext()
	bg2 := a.ChildContext(bg1)
	bg3 := a.ChildContext(bg2)

	qv, ok := bg3.Get("do")
	as.True(ok)
	if _, ok := qv.(a.Applicable); !ok {
		as.Fail("returned value not Applicable")
	}
}

func TestQuote(t *testing.T) {
	as := assert.New(t)

	r1 := runCode("(quote (blah 2 3))").(*a.List)
	r2 := runCode("'(blah 2 3)").(*a.List)

	v1, ok := r1.Get(0)
	v2, _ := r2.Get(0)
	as.True(ok, "first element retrieved")
	as.Equal(v1, v2, "first element same")

	v1, ok = r1.Get(0)
	as.True(ok, "first element retrieved")
	if _, ok := v1.(a.Symbol); !ok {
		as.Fail("first element is not a symbol")
	}

	v1, ok = r1.Get(1)
	v2, _ = r2.Get(1)
	as.True(ok, "first element retrieved")
	as.Equal(v1, v2, "second element same")

	v1, ok = r1.Get(1)
	as.Equal(a.EqualTo, a.NewFloat(2.0).Cmp(v1.(*a.Number)), "second")

	v1, ok = r1.Get(2)
	v2, _ = r2.Get(2)
	as.Equal(v1, v2, "third element same")

	v1, ok = r1.Get(2)
	as.True(ok, "third element retrieved")
	as.Equal(a.EqualTo, a.NewFloat(3.0).Cmp(v1.(*a.Number)), "third")
}

func TestDo(t *testing.T) {
	testCode(t, `
		(do
			55
		    (pr "hello")
			(prn "there" 45)
			(print "how ")
			(println "are" 66)
			(if true 99 33))
	`, a.NewFloat(99))
}

func TestTrueFalse(t *testing.T) {
	testCode(t, `true`, a.True)
	testCode(t, `false`, a.False)
	testCode(t, `nil`, a.Nil)
}
