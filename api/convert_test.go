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

	as.String(`["hello" "there"]`, v1)
	as.Identical(v1, v2)
	as.String(`("hello" "there")`, l2)
	as.Identical(l2, l3)
	as.String(`{"hello" "there"}`, a1)
	as.Identical(a1, a2)
}
