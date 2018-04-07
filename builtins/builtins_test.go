package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	b "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

func getBuiltIn(n a.Name) a.Function {
	if r, ok := b.GetFunction(n); ok {
		return r
	}
	panic(a.ErrStr("Built in not found: ", n))
}

func f(n float64) a.Number {
	return a.NewFloat(n)
}

func v(e ...a.Value) a.Vector {
	return a.Vector(e)
}

func args(e ...a.Value) a.Vector {
	return a.Vector(e)
}

func s(v string) a.Str {
	return a.Str(v)
}

func kw(n a.Name) a.Keyword {
	return a.NewKeyword(n)
}

func local(n a.Name) a.Symbol {
	return a.NewLocalSymbol(n)
}

func TestMissingBuiltIn(t *testing.T) {
	as := assert.New(t)

	r, ok := b.GetFunction("boom")
	as.Nil(r)
	as.False(ok)
}

func TestExplodingBuiltInCall(t *testing.T) {
	as := assert.New(t)
	read := getBuiltIn("read")
	eval := getBuiltIn("eval")

	defer as.ExpectError(a.ErrStr(a.KeyNotFound, a.Name("boom")))
	c := e.NewEvalContext()
	r := read.Apply(c, args(s("(def-builtin boom)")))
	eval.Apply(c, args(r))
}
