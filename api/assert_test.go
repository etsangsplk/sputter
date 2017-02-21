package api_test

import (
	"fmt"
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

func TestBadArity(t *testing.T) {
	a := assert.New(t)
	v := &s.Vector{1, 2, 3}

	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Sprintf(s.BadArity, 4, 3)
			a.Equal(err, rec, "arity error properly raised")
			return
		}
		a.Fail("arity error not raised")
	}()

	s.AssertArity(v, 4)
}

func TestMinimumArity(t *testing.T) {
	a := assert.New(t)
	v := &s.Vector{1, 2, 3}

	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Sprintf(s.BadMinimumArity, 4, 3)
			a.Equal(err, rec, "arity error properly raised")
			return
		}
		a.Fail("arity error not raised")
	}()

	s.AssertMinimumArity(v, 4)
}

func TestArityRange(t *testing.T) {
	a := assert.New(t)
	v := &s.Vector{1, 2, 3}

	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Sprintf(s.BadArityRange, 4, 7, 3)
			a.Equal(err, rec, "arity error properly raised")
			return
		}
		a.Fail("arity error not raised")
	}()

	s.AssertArityRange(v, 4, 7)
}
