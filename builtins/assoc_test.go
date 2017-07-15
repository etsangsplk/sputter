package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	e "github.com/kode4food/sputter/evaluator"
)

func TestAssoc(t *testing.T) {
	as := assert.New(t)
	c := e.NewEvalContext()

	assoc := getBuiltIn("assoc")
	a1 := assoc(c, args(kw("hello"), s("foo")))
	m1 := a.AssertMappedSequence(a1)
	v1, ok := m1.Get(kw("hello"))
	as.True(ok)
	as.String("foo", v1)

	isAssoc := getBuiltIn("assoc?")
	as.True(isAssoc(c, args(a1)))
	as.False(isAssoc(c, args(f(99))))

	isMapped := getBuiltIn("mapped?")
	as.True(isMapped(c, args(a1)))
	as.False(isMapped(c, args(f(99))))
}
