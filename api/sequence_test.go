package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	as := assert.New(t)

	add := a.NewFunction(
		func(_ a.Context, args a.Sequence) a.Value {
			v := args.(a.Vector)
			return v[0].(int) + v[1].(int)
		},
	)

	as.Equal(30, a.Reduce(nil, a.Vector{10, 20}, add))
	as.Equal(60, a.Reduce(nil, a.Vector{10, 20, 30}, add))
	as.Equal(100, a.Reduce(nil, a.Vector{10, 20, 30, 40}, add))

	err := a.Err(a.BadMinimumArity, 2, 1)
	defer expectError(as, err)
	a.Reduce(nil, a.Vector{10}, add)
}

type noCountSequence struct{}

func (n *noCountSequence) First() a.Value               { return nil }
func (n *noCountSequence) Rest() a.Sequence             { return nil }
func (n *noCountSequence) Prepend(v a.Value) a.Sequence { return nil }
func (n *noCountSequence) IsSequence() bool             { return true }

func TestNonCountableSequence(t *testing.T) {
	as := assert.New(t)
	nc := &noCountSequence{}

	defer expectError(as, a.ExpectedCountable)
	a.Count(nc)
}

func TestAssertSequence(t *testing.T) {
	as := assert.New(t)
	a.AssertSequence(a.NewList("hello"))

	defer expectError(as, a.Err(a.ExpectedSequence, "99"))
	a.AssertSequence(a.NewFloat(99))
}

func TestAssertIndexed(t *testing.T) {
	as := assert.New(t)
	a.AssertIndexed(a.NewList("hello"))

	defer expectError(as, a.Err(a.ExpectedIndexed, "99"))
	a.AssertIndexed(a.NewFloat(99))
}

func TestAssertConjoiner(t *testing.T) {
	as := assert.New(t)
	a.AssertConjoiner(a.NewList("hello"))

	defer expectError(as, a.Err(a.ExpectedConjoiner, "99"))
	a.AssertConjoiner(a.NewFloat(99))
}
