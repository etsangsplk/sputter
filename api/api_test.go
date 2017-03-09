package api_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestTruthy(t *testing.T) {
	as := assert.New(t)

	as.True(a.Truthy(a.True), "API True is Truthy")
	as.True(a.Truthy(true), "true is Truthy")
	as.True(a.Truthy(a.NewList("Hello")), "Non-Empty List Is Truthy")
	as.True(a.Truthy("hello"), "String is Truthy")

	as.False(a.Truthy(a.Nil), "API Nil is not Truthy")
	as.False(a.Truthy(nil), "nil is not Truthy")
	as.False(a.Truthy(a.False), "API False is not Truthy")
	as.False(a.Truthy(false), "false is not Truthy")
}

type testSequence struct{}

func (t *testSequence) Iterate() a.Iterator {
	return nil
}

func TestNonFiniteCount(t *testing.T) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.Equal(a.ExpectedFinite, rec, "count panics properly")
			return
		}
		as.Fail("count should panic")
	}()

	i := &testSequence{}
	a.Count(i)
}

func TestAssertSequence(t *testing.T) {
	as := assert.New(t)
	a.AssertSequence(a.NewList("hello"))

	defer expectError(as, a.ExpectedSequence)
	a.AssertSequence(big.NewFloat(99))
}

func TestAssertNumeric(t *testing.T) {
	as := assert.New(t)
	a.AssertNumeric(big.NewFloat(99))

	defer expectError(as, a.ExpectedNumeric)
	a.AssertNumeric(&a.Symbol{})
}
