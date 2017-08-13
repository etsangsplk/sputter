package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

type ncSeq struct{}

func (n *ncSeq) First() a.Value                     { return nil }
func (n *ncSeq) Rest() a.Sequence                   { return nil }
func (n *ncSeq) Split() (a.Value, a.Sequence, bool) { return nil, nil, false }
func (n *ncSeq) Prepend(v a.Value) a.Sequence       { return nil }
func (n *ncSeq) IsSequence() bool                   { return false }
func (n *ncSeq) Str() a.Str                         { return s("()") }

func TestNonCountableSequence(t *testing.T) {
	as := assert.New(t)
	nc := new(ncSeq)

	defer as.ExpectError(a.ErrStr(a.ExpectedCounted, "()"))
	a.Count(nc)
}

func TestAssertSequence(t *testing.T) {
	as := assert.New(t)
	a.AssertSequence(a.NewList(s("hello")))

	defer as.ExpectError(a.ErrStr(a.ExpectedSequence, f(99)))
	a.AssertSequence(f(99))
}

func TestAssertIndexed(t *testing.T) {
	as := assert.New(t)
	a.AssertIndexed(a.NewList(s("hello")))

	defer as.ExpectError(a.ErrStr(a.ExpectedIndexed, f(99)))
	a.AssertIndexed(f(99))
}

func TestAssertConjoiner(t *testing.T) {
	as := assert.New(t)
	a.AssertConjoiner(a.NewList(s("hello")))

	defer as.ExpectError(a.ErrStr(a.ExpectedConjoiner, f(99)))
	a.AssertConjoiner(f(99))
}
