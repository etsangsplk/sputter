package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

var helloName = &s.Function{
	Name: "hello",
	Exec: func(c s.Context, args s.Sequence) s.Value {
		i := args.Iterate()
		a, _ := i.Next()
		v := evaluateToString(c, a)
		return "Hello, " + v + "!"
	},
}

func evaluateToString(c s.Context, v s.Value) string {
	return s.String(s.Eval(c, v))
}

func TestEvaluate(t *testing.T) {
	a := assert.New(t)

	l := &s.Cons{Car: helloName, Cdr: s.NewList("World")}
	r := s.Eval(s.NewContext(), l)

	a.Equal("Hello, World!", r.(string), "good hello")
}

func TestEvaluateSequence(t *testing.T) {
	a := assert.New(t)

	s1 := &s.Cons{Car: helloName, Cdr: s.NewList("World")}
	s2 := &s.Cons{Car: helloName, Cdr: s.NewList("Foo")}
	l := &s.Cons{Car: s1, Cdr: s.NewList(s2)}

	r := s.EvalSequence(s.NewContext(), l)
	a.Equal("Hello, Foo!", r.(string), "last result")
}
