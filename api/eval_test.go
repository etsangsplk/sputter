package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

var helloName = &a.Function{
	Name: "hello",
	Apply: func(c a.Context, args a.Sequence) a.Value {
		i := args.Iterate()
		a, _ := i.Next()
		v := evaluateToString(c, a)
		return "Hello, " + v + "!"
	},
}

func evaluateToString(c a.Context, v a.Value) string {
	return a.String(a.Eval(c, v))
}

func TestEvaluate(t *testing.T) {
	as := assert.New(t)

	l := &a.Cons{Car: helloName, Cdr: a.NewList("World")}
	c := a.NewContext()
	r := a.Eval(c, l)

	as.Equal("Hello, World!", r.(string), "good hello")
}

func TestEvaluateSequence(t *testing.T) {
	as := assert.New(t)

	s1 := &a.Cons{Car: helloName, Cdr: a.NewList("World")}
	s2 := &a.Cons{Car: helloName, Cdr: a.NewList("Foo")}
	l := &a.Cons{Car: s1, Cdr: a.NewList(s2)}

	c := a.NewContext()
	r := a.EvalSequence(c, l)
	as.Equal("Hello, Foo!", r.(string), "last result")
}
