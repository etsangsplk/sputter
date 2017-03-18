package api_test

import (
	"fmt"
	"strings"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {
	as := assert.New(t)

	f1 := &a.Function{
		Name: "test-function",
		Doc:  "this is a test function",
		Exec: func(c a.Context, args a.Sequence) a.Value {
			return "hello"
		},
	}

	f2 := &a.Function{}

	as.True(strings.HasPrefix(f1.String(), "(fn :name"), "name returned")
	as.Equal("this is a test function", f1.Docstring(), "doc returned")
	as.True(strings.HasPrefix(f2.String(), "(fn :addr"), "address returned")

	c := a.NewContext()
	as.Equal("hello", f1.Apply(c, a.EmptyList), "function executes")
}

func TestGoodArity(t *testing.T) {
	as := assert.New(t)
	v := &a.Vector{1, 2, 3}

	defer func() {
		if rec := recover(); rec != nil {
			as.Fail("arity tests should not explode")
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

func TestAssertApplicable(t *testing.T) {
	as := assert.New(t)
	a.AssertApplicable(&a.Function{})

	defer expectError(as, a.ExpectedApplicable)
	a.AssertApplicable(&a.Symbol{})
}
