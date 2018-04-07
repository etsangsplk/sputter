package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	e "github.com/kode4food/sputter/evaluator"
)

func TestClosure(t *testing.T) {
	as := assert.New(t)
	makeClosure := getBuiltIn("make-closure")

	c := e.NewEvalContext()
	r1 := makeClosure.Apply(c, args(
		v(local("ignore")),
		v(s("hello"), local("name"), local("ignore")),
	))

	as.String(`(sputter:closure [name] [["hello" name ignore]])`, r1)

	closure := getBuiltIn("closure")
	c.Put("ignore", f(99))
	c.Put("name", s("Bob"))
	ar := r1.(*a.List).Rest()

	defer as.ExpectError(a.ErrStr(a.UnknownSymbol, a.Name("ignore")))
	a.Apply(c, closure, ar)
}
