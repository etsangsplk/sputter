package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestTruthy(t *testing.T) {
	a := assert.New(t)

	a.True(s.Truthy(s.True), "API True is Truthy")
	a.True(s.Truthy(true), "true is Truthy")
	a.True(s.Truthy(s.NewList("Hello")), "Non-Empty List Is Truthy")
	a.True(s.Truthy("hello"), "String is Truthy")

	a.False(s.Truthy(s.Nil), "API Nil is not Truthy")
	a.False(s.Truthy(nil), "nil is not Truthy")
	a.False(s.Truthy(s.False), "API False is not Truthy")
	a.False(s.Truthy(false), "false is not Truthy")
}

type testSequence struct{}

func (t *testSequence) Iterate() s.Iterator {
	return nil
}

func TestNonFiniteCount(t *testing.T) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(s.NonFinite, rec, "count panics properly")
			return
		}
		a.Fail("count should panic")
	}()

	i := &testSequence{}
	s.Count(i)
}
