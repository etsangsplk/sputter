package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestVector(t *testing.T) {
	as := assert.New(t)

	v1 := &a.Vector{"hello", "how", "are", "you?"}
	as.Equal(4, v1.Count(), "vector 1 count is correct")
	as.Equal(4, a.Count(v1), "vector 1 general count is correct")
	as.Equal("are", v1.Get(2), "get by index is correct")
	as.Equal(`["hello" "how" "are" "you?"]`, v1.String(), "string is good")

	v2 := v1.Prepend("oh").(a.Vector)
	as.Equal(5, v2.Count(), "vector 2 count is correct")
	as.Equal(4, v1.Count(), "vector 1 count is still correct")
	as.Equal("oh", v2.Get(0), "get by index is correct")
	as.Equal("are", v2.Get(3), "get by index is correct")

	c := a.NewContext()
	as.Equal("are", v1.Apply(c, a.NewList(a.NewFloat(2))))
}

type testEvaluable struct{}

func (t *testEvaluable) Eval(c a.Context) a.Value {
	return "are"
}

func TestVectorEval(t *testing.T) {
	as := assert.New(t)

	v := &a.Vector{"hello", "how", &testEvaluable{}, "you?"}
	c := a.NewContext()
	r := v.Eval(c)

	if _, ok := r.(a.Indexed); !ok {
		as.Fail("result is not a finite sequence")
	}

	as.Equal("are", r.(a.Indexed).Get(2), "get is working")
	as.Equal(`["hello" "how" "are" "you?"]`, a.String(r), "string is good")
}

func TestIterate(t *testing.T) {
	as := assert.New(t)

	v := &a.Vector{"hello", "how", "are", "you?"}
	i := a.Iterate(v)
	e1, _ := i.Next()
	s1 := i.Rest()
	e2, _ := i.Next()
	s2 := i.Rest()
	e3, _ := i.Next()
	e4, _ := i.Next()
	e5, ok := i.Next()

	as.Equal("hello", e1, "first vector element")
	as.Equal("how", e2, "second vector element")
	as.Equal("are", e3, "third vector element")
	as.Equal("you?", e4, "fourth vector element")

	as.Equal(3, a.Count(s1), "s1 slice count")
	as.Equal(2, a.Count(s2), "s2 slice count")

	as.Equal(a.Nil, e5, "fifth element is nil")
	as.False(ok, "fifth element was false")
}
