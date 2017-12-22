package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	e "github.com/kode4food/sputter/evaluator"
)

func TestClosure(t *testing.T) {
	as := assert.New(t)

	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("chain")
	ns.Put("chain", s("hello"))

	c1 := a.ChildContext(ns, a.Variables{})
	c1.Put("foo", s("foo_val"))
	c1.Put("bar", s("bar_val"))
	c1.Put("baz", s("baz_val"))

	cl := e.EvalStr(c1, "(let [p (promise)] (p [foo, baz]) (p))")
	v1 := a.Eval(c1, cl).(a.Vector)

	r1, _ := v1.ElementAt(0)
	r2, _ := v1.ElementAt(1)
	as.String("foo_val", r1)
	as.String("baz_val", r2)
}
