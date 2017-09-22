package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestVector(t *testing.T) {
	as := assert.New(t)

	v1 := a.NewVector(s("hello"), s("how"), s("are"), s("you?"))
	as.Number(4, v1.Count())
	as.Number(4, a.Count(v1))

	r, ok := v1.ElementAt(2)
	as.True(ok)
	as.String("are", r)
	as.String(`["hello" "how" "are" "you?"]`, v1)

	v2 := v1.Prepend(s("oh")).(a.Vector)
	as.Number(5, v2.Count())
	as.Number(4, v1.Count())

	v3 := v2.Conjoin(s("good?")).(a.Vector)
	r, ok = v3.ElementAt(5)
	as.True(ok)
	as.String("good?", r)
	as.Number(6, v3.Count())

	r, ok = v3.ElementAt(0)
	as.True(ok)
	as.String("oh", r)

	r, ok = v3.ElementAt(3)
	as.True(ok)
	as.String("are", r)

	c := a.NewContext()
	as.String("are", v1.Apply(c, a.NewList(f(2))))
}

func TestEmptyVector(t *testing.T) {
	as := assert.New(t)

	v := a.NewVector()
	as.Nil(v.First())
	as.String("[]", v.Str())
	as.String("[]", v.Rest())
}

type testEvaluable struct{}

func (t *testEvaluable) Eval(_ a.Context) a.Value {
	return s("are")
}

func (t *testEvaluable) Str() a.Str {
	return s("")
}

func TestVectorEval(t *testing.T) {
	as := assert.New(t)

	v := a.NewVector(s("hello"), s("how"), new(testEvaluable), s("you?"))

	c := a.NewContext()
	r := a.Eval(c, v)

	if _, ok := r.(a.Indexed); !ok {
		as.Fail("result is not a finite sequence")
	}

	i, ok := r.(a.Indexed).ElementAt(2)
	as.True(ok)
	as.String("are", i)
	as.String(`["hello" "how" "are" "you?"]`, r)
}

func TestIterate(t *testing.T) {
	as := assert.New(t)

	v := a.NewVector(s("hello"), s("how"), s("are"), s("you?"))
	i := a.Iterate(v)
	e1, _ := i.Next()
	s1 := i.Rest()
	e2, _ := i.Next()
	s2 := i.Rest()
	e3, _ := i.Next()
	e4, _ := i.Next()
	e5, ok := i.Next()

	as.String("hello", e1)
	as.String("how", e2)
	as.String("are", e3)
	as.String("you?", e4)

	as.Number(3, a.Count(s1))
	as.Number(2, a.Count(s2))

	as.Equal(a.Nil, e5)
	as.False(ok)
}

func TestVectorExplosion(t *testing.T) {
	as := assert.New(t)

	idx := f(3)
	err := a.ErrStr(a.IndexNotFound, idx)
	defer as.ExpectError(err)

	v := a.NewVector(s("foo"))
	v.Apply(a.NewContext(), a.NewList(idx))
}
