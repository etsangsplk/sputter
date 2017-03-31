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

	f1 := a.NewFunction(func(_ a.Context, _ a.Sequence) a.Value {
		return "hello"
	}).WithMetadata(a.Metadata{
		a.MetaName: a.Name("test-function"),
		a.MetaDoc:  "this is a test",
	}).(a.Function)

	f2 := a.NewFunction(nil)
	f3 := f1.WithMetadata(a.Metadata{a.MetaDoc: "modified"})

	as.NotNil(f1.Metadata())
	as.NotNil(f2.Metadata())
	as.NotEqual(f1.Metadata(), f3.Metadata())

	as.Equal("this is a test", f1.Metadata()[a.MetaDoc], "not modified")
	as.Equal("modified", f3.Metadata()[a.MetaDoc], "modified")

	as.True(strings.Contains(a.String(f1), ":name test-function"), "name")
	as.Equal("this is a test", f1.Documentation(), "doc returned")
	as.True(strings.Contains(a.String(f2), ":name <lambda>"), "lambda")

	c := a.NewContext()
	as.Equal("hello", f1.Apply(c, a.EmptyList), "function executes")
}

func TestGoodArity(t *testing.T) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.Fail("arity tests should not explode")
		}
	}()

	v1 := &a.Vector{1, 2, 3}
	as.Equal(3, a.AssertArity(v1, 3))
	as.Equal(3, a.AssertArityRange(v1, 2, 4))
	as.Equal(3, a.AssertArityRange(v1, 3, 3))
	as.Equal(3, a.AssertMinimumArity(v1, 3))
	as.Equal(3, a.AssertMinimumArity(v1, 2))

	v2 := a.NewConcat(a.Vector{
		a.NewList(1),
		a.Vector{2, 3},
	})
	as.Equal(3, a.AssertArity(v2, 3))
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
