package api_test

import (
	"fmt"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestFunction(t *testing.T) {
	as := assert.New(t)

	f1 := a.NewFunction(func(_ a.Context, _ a.Sequence) a.Value {
		return s("hello")
	}).WithMetadata(a.Metadata{
		a.MetaName: a.Name("test-function"),
		a.MetaDoc:  s("this is a test"),
	}).(a.Function)

	as.True(f1.IsFunction())

	f2 := a.NewFunction(nil)
	f3 := f1.WithMetadata(a.Metadata{a.MetaDoc: s("modified")})

	as.NotNil(f1.Metadata())
	as.NotNil(f2.Metadata())
	as.NotIdentical(f1.Metadata(), f3.Metadata())

	as.String("this is a test", f1.Metadata()[a.MetaDoc])
	as.String("modified", f3.Metadata()[a.MetaDoc])

	as.Contains(":name test-function", f1)
	as.String("this is a test", f1.Documentation())

	c := a.NewContext()
	as.String("hello", f1.Apply(c, a.EmptyList))

	f4 := a.NewFunction(nil).WithMetadata(a.Metadata{
		a.MetaType: f(99),
	}).(a.Function)

	as.String("function", f4.Type())
}

func TestGoodArity(t *testing.T) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.Fail("arity tests should not explode")
		}
	}()

	v1 := a.NewVector(f(1), f(2), f(3))
	as.Number(3, a.AssertArity(v1, 3))
	as.Number(3, a.AssertArityRange(v1, 2, 4))
	as.Number(3, a.AssertArityRange(v1, 3, 3))
	as.Number(3, a.AssertMinimumArity(v1, 3))
	as.Number(3, a.AssertMinimumArity(v1, 2))

	v2 := a.Concat(a.NewVector(
		a.NewList(f(1)),
		a.NewVector(f(2), f(3)),
	))
	as.Number(3, a.AssertArity(v2, 3))
}

func TestBadArity(t *testing.T) {
	as := assert.New(t)
	v := a.NewVector(f(1), f(2), f(3))

	defer as.ExpectError(a.Err(a.BadArity, 4, 3))
	a.AssertArity(v, 4)
}

func TestMinimumArity(t *testing.T) {
	as := assert.New(t)
	v := a.NewVector(f(1), f(2), f(3))

	defer as.ExpectError(a.Err(a.BadMinimumArity, 4, 3))
	a.AssertMinimumArity(v, 4)
}

func TestArityRange(t *testing.T) {
	as := assert.New(t)
	v := a.NewVector(f(1), f(2), f(3))

	defer as.ExpectError(fmt.Sprintf(a.BadArityRange, 4, 7, 3))
	a.AssertArityRange(v, 4, 7)
}
