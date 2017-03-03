package api_test

import (
	"fmt"
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestGoodArity(t *testing.T) {
	a := assert.New(t)
	v := &s.Vector{1, 2, 3}

	defer func() {
		if rec := recover(); rec != nil {
			a.Fail("arity tests should not explode")
			return
		}
	}()

	s.AssertArity(v, 3)
	s.AssertArityRange(v, 2, 4)
	s.AssertArityRange(v, 3, 3)
	s.AssertMinimumArity(v, 3)
	s.AssertMinimumArity(v, 2)
}

func expectError(a *assert.Assertions, err string) {
	if rec := recover(); rec != nil {
		a.Equal(err, rec, "error raised")
		return
	}
	a.Fail("error not raised")
}

func TestBadArity(t *testing.T) {
	a := assert.New(t)
	v := &s.Vector{1, 2, 3}

	defer expectError(a, fmt.Sprintf(s.BadArity, 4, 3))
	s.AssertArity(v, 4)
}

func TestMinimumArity(t *testing.T) {
	a := assert.New(t)
	v := &s.Vector{1, 2, 3}

	defer expectError(a, fmt.Sprintf(s.BadMinimumArity, 4, 3))
	s.AssertMinimumArity(v, 4)
}

func TestArityRange(t *testing.T) {
	a := assert.New(t)
	v := &s.Vector{1, 2, 3}

	defer expectError(a, fmt.Sprintf(s.BadArityRange, 4, 7, 3))
	s.AssertArityRange(v, 4, 7)
}

func TestAssertCons(t *testing.T) {
	a := assert.New(t)
	s.AssertCons(&s.Cons{Car: "hello", Cdr: "there"})

	defer expectError(a, s.ExpectedCons)
	s.AssertCons(big.NewFloat(99))
}

func TestAssertSequence(t *testing.T) {
	a := assert.New(t)
	s.AssertSequence(s.NewList("hello"))

	defer expectError(a, s.ExpectedSequence)
	s.AssertSequence(big.NewFloat(99))
}

func TestAssertSymbol(t *testing.T) {
	a := assert.New(t)
	s.AssertSymbol(&s.Symbol{})

	defer expectError(a, s.ExpectedSymbol)
	s.AssertSymbol(big.NewFloat(99))
}

func TestAssertNumeric(t *testing.T) {
	a := assert.New(t)
	s.AssertNumeric(big.NewFloat(99))

	defer expectError(a, s.ExpectedNumeric)
	s.AssertNumeric(&s.Symbol{})
}
 
func TestAssertFunction(t *testing.T) {
	a := assert.New(t)
	s.AssertFunction(&s.Function{})

	defer expectError(a, s.ExpectedFunction)
	s.AssertFunction(&s.Symbol{})
}
