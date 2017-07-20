package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestSequenceConversions(t *testing.T) {
	as := assert.New(t)
	l1 := a.NewList(a.Str("hello"), a.Str("there"))
	v1 := a.ToVector(l1)
	v2 := a.ToVector(v1)
	l2 := a.ToList(v2)
	l3 := a.ToList(l2)
	a1 := a.ToAssociative(l3)
	a2 := a.ToAssociative(a1)

	l4 := a.NewList(a.Str("hello"), a.Nil, a.Str("there"), v1)
	s1 := a.ToStr(l4)
	s2 := a.ToStr(s1)
	s3 := a.ToReaderStr(s2)

	as.String(`["hello" "there"]`, v1)
	as.Identical(v1, v2)
	as.String(`("hello" "there")`, l2)
	as.Identical(l2, l3)
	as.String(`{"hello" "there"}`, a1)
	as.Identical(a1, a2)
	as.String(`hellothere["hello" "there"]`, s1)
	as.Identical(s1, s2)
	as.Identical(s1, s3)
}

func identity(v a.Value) a.Value {
	return v
}

func TestUncountedConversions(t *testing.T) {
	as := assert.New(t)
	l1 := a.Map(a.NewList(a.Str("hello"), a.Str("there")), identity)
	v1 := a.ToVector(l1)
	v2 := a.ToVector(v1)
	l2 := a.ToList(a.Map(v2, identity))
	l3 := a.ToList(l2)
	a1 := a.ToAssociative(a.Map(l3, identity))
	a2 := a.ToAssociative(a1)

	l4 := a.Map(a.NewList(a.Str("hello"), a.Nil, a.Str("there"), v1), identity)
	s1 := a.ToStr(l4)

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
	a.ToAssociative(v1)
}

func TestUncountedAssocConvertError(t *testing.T) {
	as := assert.New(t)

	v1 := a.Map(a.NewVector(a.NewKeyword("boom")), identity)
	defer as.ExpectError(a.ErrStr(a.ExpectedPair))
	a.ToAssociative(v1)
}

func TestToReaderStr(t *testing.T) {
	as := assert.New(t)

	v1 := a.NewVector(a.NewKeyword("boom"), a.Str("hello"))
	s1 := a.ToReaderStr(v1)
	as.String(":boom \"hello\"", s1)
}
