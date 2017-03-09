package api_test

import (
	"fmt"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {
	as := assert.New(t)

	f := &a.Function{
		Name: "test-function",
		Apply: func(c a.Context, args a.Sequence) a.Value {
			return "hello"
		},
	}

	as.Equal("test-function", f.String(), "string returned")

	c := a.NewContext()
	as.Equal("hello", f.Apply(c, a.EmptyList), "function executes")
}

func TestGoodArity(t *testing.T) {
	as := assert.New(t)
	v := &a.Vector{1, 2, 3}

	defer func() {
		if rec := recover(); rec != nil {
			as.Fail("arity tests should not explode")
			return
		}
	}()

	a.AssertArity(v, 3)
	a.AssertArityRange(v, 2, 4)
	a.AssertArityRange(v, 3, 3)
	a.AssertMinimumArity(v, 3)
	a.AssertMinimumArity(v, 2)
}

func TestBadArity(t *testing.T) {
	as := assert.New(t)
	v := &a.Vector{1, 2, 3}

	defer expectError(as, fmt.Sprintf(a.BadArity, 4, 3))
	a.AssertArity(v, 4)
}

func TestMinimumArity(t *testing.T) {
	as := assert.New(t)
	v := &a.Vector{1, 2, 3}

	defer expectError(as, fmt.Sprintf(a.BadMinimumArity, 4, 3))
	a.AssertMinimumArity(v, 4)
}

func TestArityRange(t *testing.T) {
	as := assert.New(t)
	v := &a.Vector{1, 2, 3}

	defer expectError(as, fmt.Sprintf(a.BadArityRange, 4, 7, 3))
	a.AssertArityRange(v, 4, 7)
}

func TestAssertFunction(t *testing.T) {
	as := assert.New(t)
	a.AssertFunction(&a.Function{})

	defer expectError(as, a.ExpectedFunction)
	a.AssertFunction(&a.Symbol{})
}
