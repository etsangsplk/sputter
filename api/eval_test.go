package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func evaluateToString(c *s.Context, v s.Value) string {
	return s.String(s.Eval(c, v))
}

func TestEvaluate(t *testing.T) {
	a := assert.New(t)

	f := &s.Function{
		Name: "hello",
		Exec: func(c *s.Context, args s.Sequence) s.Value {
			i := args.Iterate()
			a, _ := i.Next()
			v := evaluateToString(c, a)
			return "Hello, " + v + "!"
		},
	}

	l := &s.Cons{Car: f, Cdr: s.NewList("World")}
	r := s.Eval(s.NewContext(), l)

	a.Equal("Hello, World!", r.(string), "good hello")
}
