package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestAssertSequence(t *testing.T) {
	as := assert.New(t)
	a.AssertSequence(a.NewList("hello"))

	defer expectError(as, a.ExpectedSequence)
	a.AssertSequence(a.NewFloat(99))
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
