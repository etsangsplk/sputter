package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestVector(t *testing.T) {
	a := assert.New(t)

	v := &s.Vector{"hello", "how", "are", "you?"}
	a.Equal(4, v.Count(), "vector count is correct")
	a.Equal(4, s.Count(v), "vector general count is correct")
	a.Equal("are", v.Get(2), "get by index is correct")
	a.Equal("[hello how are you?]", v.String(), "string version is good")
}

type testEvaluable struct{}

func (t *testEvaluable) Eval(c *s.Context) s.Value {
	return "are"
}

func TestVectorEval(t *testing.T) {
	a := assert.New(t)

	v := &s.Vector{"hello", "how", &testEvaluable{}, "you?"}
	r := v.Eval(s.NewContext())

	if _, ok := r.(s.Finite); !ok {
		a.Fail("result is not a finite sequence")
	}

	a.Equal("are", r.(s.Finite).Get(2), "get is working")
	a.Equal("[hello how are you?]", s.String(r), "string version is good")
}

func TestIterate(t *testing.T) {
	a := assert.New(t)

	v := &s.Vector{"hello", "how", "are", "you?"}
	i := v.Iterate()
	e1, _ := i.Next()
	s1 := i.Rest()
	e2, _ := i.Next()
	s2 := i.Rest()
	e3, _ := i.Next()
	e4, _ := i.Next()
	e5, ok := i.Next()

	a.Equal("hello", e1, "first vector element")
	a.Equal("how", e2, "second vector element")
	a.Equal("are", e3, "third vector element")
	a.Equal("you?", e4, "fourth vector element")

	a.Equal(3, s1.(s.Finite).Count(), "s1 slice count")
	a.Equal(2, s2.(s.Finite).Count(), "s2 slice count")

	a.Equal(s.Nil, e5, "fifth element is nil")
	a.False(ok, "fifth element was false")
}
