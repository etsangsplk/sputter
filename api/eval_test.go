package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

var helloName = a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
	i := a.Iterate(args)
	n, _ := i.Next()
	v := a.Eval(c, n)
	return s("Hello, " + string(v.(a.Str)) + "!")
}).WithMetadata(a.Metadata{
	a.MetaName: a.Name("hello"),
}).(*a.Function)

func TestEvaluate(t *testing.T) {
	as := assert.New(t)

	l := a.NewExpression(a.NewList(s("World")).Prepend(helloName).(*a.List))
	c := a.NewContext()
	r := a.Eval(c, l)

	as.String("Hello, World!", r)
}

func TestEvaluateSequence(t *testing.T) {
	as := assert.New(t)

	s1 := a.NewExpression(a.NewList(s("World")).Prepend(helloName).(*a.List))
	s2 := a.NewExpression(a.NewList(s("Foo")).Prepend(helloName).(*a.List))
	l := a.NewList(s2).Prepend(s1)

	c := a.NewContext()
	r := a.EvalSequence(c, l)
	as.String("Hello, Foo!", r)
}

func TestAssertApplicable(t *testing.T) {
	as := assert.New(t)
	a.AssertApplicable(a.NewFunction(nil))

	defer as.ExpectError(a.Err(a.ExpectedApplicable, f(99)))
	a.AssertApplicable(f(99))
}
