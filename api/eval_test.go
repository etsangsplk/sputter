package api_test

import (
	"fmt"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

var helloName = a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
	i := a.Iterate(args)
	n, _ := i.Next()
	v := a.Eval(c, n)
	return "Hello, " + a.String(v) + "!"
}).WithMetadata(a.Metadata{
	a.MetaName: a.Name("hello"),
})

func TestEvaluate(t *testing.T) {
	as := assert.New(t)

	l := a.NewList("World").Prepend(helloName)
	c := a.NewContext()
	r := a.Eval(c, l)

	as.Equal("Hello, World!", r.(string), "good hello")
}

func TestEvaluateSequence(t *testing.T) {
	as := assert.New(t)

	s1 := a.NewList("World").Prepend(helloName)
	s2 := a.NewList("Foo").Prepend(helloName)
	l := a.NewList(s2).Prepend(s1)

	c := a.NewContext()
	r := a.EvalSequence(c, l)
	as.Equal("Hello, Foo!", r.(string), "last result")
}

func TestAssertApplicable(t *testing.T) {
	as := assert.New(t)
	a.AssertApplicable(a.NewFunction(nil))

	defer expectError(as, a.ExpectedApplicable)
	a.AssertApplicable(a.NewFloat(99))
}
