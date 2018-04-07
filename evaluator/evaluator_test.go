package evaluator_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	e "github.com/kode4food/sputter/evaluator"
	r "github.com/kode4food/sputter/reader"
)

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) a.Number {
	return a.NewFloat(f)
}

func TestEvalContext(t *testing.T) {
	as := assert.New(t)

	uc := a.GetNamespace(a.UserDomain)
	uc.Delete("foo")
	uc.Put("foo", f(99))

	ec := e.NewEvalContext()
	v, _ := ec.Get("foo")
	as.Number(99, v)
}

func testCodeWithContext(
	as *assert.Wrapper, code string, expect a.Value, c a.Context) {
	as.Equal(expect, e.EvalStr(c, s(code)))
}

func TestEvaluable(t *testing.T) {
	as := assert.New(t)

	hello := a.NewExecFunction(func(c a.Context, args a.Vector) a.Value {
		i := a.Iterate(args)
		arg, _ := i.Next()
		v := a.Eval(c, arg)
		return s("Hello, " + string(v.(a.Str)) + "!")
	}).WithMetadata(a.Properties{
		a.NameKey: a.Name("hello"),
	}).(a.Function)

	c := e.NewEvalContext()
	c.Put("hello", hello)
	c.Put("name", s("Bob"))

	testCodeWithContext(as, `(hello "World")`, s("Hello, World!"), c)
	testCodeWithContext(as, `(hello name)`, s("Hello, Bob!"), c)
}

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)

	b := e.NewEvalContext()
	ns := a.GetContextNamespace(b)
	ns.Delete("hello")

	ns.Put("hello", a.NewExecFunction(func(_ a.Context, _ a.Vector) a.Value {
		return s("there")
	}).WithMetadata(a.Properties{
		a.NameKey: a.Name("hello"),
	}).(a.Function))

	l := r.Scan(`(hello)`)
	tr := r.Read(l)

	c := a.ChildLocals(b)
	ev := e.Evaluate(c, tr)
	v, _ := a.Last(ev)
	as.String("there", v)
}
