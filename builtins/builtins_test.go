package builtins_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	r "github.com/kode4food/sputter/reader"
	"github.com/stretchr/testify/assert"
)

func runCode(src string) s.Value {
	l := r.NewLexer(src)
	tr := r.NewReader(b.Context, l)
	return r.EvalReader(s.NewContext(), tr)
}

func testCode(t *testing.T, src string, expect s.Value) {
	a := assert.New(t)
	a.Equal(expect, runCode(src), src)
}

func testBadCode(t *testing.T, src string, err string) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(err, rec, "bad code panics properly")
			return
		}
		a.Fail("bad code should panic")
	}()

	runCode(src)
}

func TestBuiltInsContext(t *testing.T) {
	a := assert.New(t)

	bg1 := s.ChildContext(b.Context)
	bg2 := s.ChildContext(bg1)
	bg3 := s.ChildContext(bg2)

	a.Equal(b.Context, bg3.Globals())

	qv, ok := bg3.Get("quote")
	a.True(ok)
	if fv, ok := qv.(*s.Function); ok {
		a.Equal("quote", string(fv.Name))
	} else {
		a.Fail("returned value not a Function")
	}
}

func TestQuote(t *testing.T) {
	a := assert.New(t)

	r1 := runCode("(quote (1 2 3))").(*s.Cons)
	r2 := runCode("'(1 2 3)").(*s.Cons)

	a.Equal(r1.Get(0), r2.Get(0), "first element same")
	a.Equal(big.NewFloat(1.0), r1.Get(0), "first element correct")
	a.Equal(r1.Get(1), r2.Get(1), "second element same")
	a.Equal(big.NewFloat(2.0), r1.Get(1), "second element correct")
	a.Equal(r1.Get(2), r2.Get(2), "third element same")
	a.Equal(big.NewFloat(3.0), r1.Get(2), "third element correct")
}

func TestDo(t *testing.T) {
	testCode(t, `
		(do
			(println "hello")
			(println "there")
			99)
	`, big.NewFloat(99))
}

func TestTrueFalse(t *testing.T) {
	testCode(t, `true`, s.True)
	testCode(t, `false`, s.False)
	testCode(t, `nil`, s.Nil)
}
