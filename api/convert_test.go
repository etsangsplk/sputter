package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestSequenceConversions(t *testing.T) {
	as := assert.New(t)
	l1 := a.NewList(a.Str("hello"), a.Str("there"))
	v1 := a.SequenceToVector(l1)
	v2 := a.SequenceToVector(v1)
	l2 := a.SequenceToList(v2)
	l3 := a.SequenceToList(l2)
	a1 := a.SequenceToAssociative(l3)
	a2 := a.SequenceToAssociative(a1)

	l4 := a.NewList(a.Str("hello"), a.Nil, a.Str("there"), v1)
	s1 := a.SequenceToStr(l4)
	s2 := a.SequenceToStr(s1)

	as.String(`["hello" "there"]`, v1)
	as.Identical(v1, v2)
	as.String(`("hello" "there")`, l2)
	as.Identical(l2, l3)
	as.String(`{"hello" "there"}`, a1)
	as.Identical(a1, a2)
	as.String(`hellothere["hello" "there"]`, s1)
	as.Identical(s1, s2)
}

func identity(v a.Value) a.Value {
	return v
}

func TestUncountedConversions(t *testing.T) {
	as := assert.New(t)
	l1 := a.Map(a.NewList(a.Str("hello"), a.Str("there")), identity)
	v1 := a.SequenceToVector(l1)
	v2 := a.SequenceToVector(v1)
	l2 := a.SequenceToList(a.Map(v2, identity))
	l3 := a.SequenceToList(l2)
	a1 := a.SequenceToAssociative(a.Map(l3, identity))
	a2 := a.SequenceToAssociative(a1)

	l4 := a.Map(a.NewList(a.Str("hello"), a.Nil, a.Str("there"), v1), identity)
	s1 := a.SequenceToStr(l4)

	as.String(`["hello" "there"]`, v1)
	as.Identical(v1, v2)
	as.String(`("hello" "there")`, l2)
	as.Identical(l2, l3)
	as.String(`{"hello" "there"}`, a1)
	as.Identical(a1, a2)
	as.String(`hellothere["hello" "there"]`, s1)
}

func TestAssocConvertError(t *testing.T) {
	as := assert.New(t)

	v1 := a.NewVector(a.NewKeyword("boom"))
	defer as.ExpectError(a.ErrStr(a.ExpectedPair))
	a.SequenceToAssociative(v1)
}

func TestUncountedAssocConvertError(t *testing.T) {
	as := assert.New(t)

	v1 := a.Map(a.NewVector(a.NewKeyword("boom")), identity)
	defer as.ExpectError(a.ErrStr(a.ExpectedPair))
	a.SequenceToAssociative(v1)
}
