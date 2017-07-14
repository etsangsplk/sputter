package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	b "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

func getBuiltIn(n a.Name) a.SequenceProcessor {
	if r, ok := b.GetBuiltIn(n); ok {
		return r
	}
	panic(a.Err("Built in not found: ", n))
}

func TestAssoc(t *testing.T) {
	as := assert.New(t)
	c := e.NewEvalContext()

	assoc := getBuiltIn("assoc")
	a1 := assoc(c, a.NewVector(a.NewKeyword("hello"), a.Str("foo")))
	m1 := a.AssertMappedSequence(a1)
	v1, ok := m1.Get(a.NewKeyword("hello"))
	as.True(ok)
	as.String("foo", v1)
}
