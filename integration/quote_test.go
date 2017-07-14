package integration_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	_ "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

func TestQuote(t *testing.T) {
	as := assert.New(t)

	r1 := runCode("(quote (blah 2 3))").(a.List)
	r2 := runCode("'(blah 2 3)").(a.List)

	v1, ok := r1.ElementAt(0)
	v2, _ := r2.ElementAt(0)
	as.True(ok)
	as.Equal(v1, v2)

	v1, ok = r1.ElementAt(0)
	as.True(ok)
	if _, ok := v1.(a.Symbol); !ok {
		as.Fail("first element is not a symbol")
	}

	v1, ok = r1.ElementAt(1)
	v2, _ = r2.ElementAt(1)
	as.True(ok)
	as.Identical(v1, v2)

	v1, ok = r1.ElementAt(1)
	as.True(ok)
	as.Number(2, v1)

	v1, ok = r1.ElementAt(2)
	v2, _ = r2.ElementAt(2)
	as.True(ok)
	as.Identical(v1, v2)

	v1, ok = r1.ElementAt(2)
	as.True(ok)
	as.Number(3, v1)
}

func TestUnquote(t *testing.T) {
	as := assert.New(t)

	c := e.NewEvalContext()
	c.Put("foo", a.NewFloat(456))
	r1 := e.EvalStr(c, `'[123 ~foo]`)
	as.String("[123 (sputter:unquote foo)]", r1)
}

func TestUnquoteMacro(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("test")

	runCode("(defmacro test [x & y] `(~x ~@y {:hello 99}))")
	testCode(t, "(test vector 1 2 3)", s("[1 2 3 {:hello 99}]"))
}
