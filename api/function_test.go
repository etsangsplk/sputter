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
	}).WithMetadata(a.Properties{
		a.MetaName: a.Name("test-function"),
		a.MetaDoc:  s("this is a test"),
	}).(a.Function)

	as.True(f1.IsFunction())

	f2 := a.NewFunction(nil)
	f3 := f1.WithMetadata(a.Properties{
		a.MetaDoc: s("modified"),
	})

	as.NotNil(f1.Metadata())
	as.NotNil(f2.Metadata())
	as.NotIdentical(f1.Metadata(), f3.Metadata())

	v, _ := f1.Metadata().Get(a.MetaDoc)
	as.String("this is a test", v)

	v, _ = f3.Metadata().Get(a.MetaDoc)
	as.String("modified", v)

	as.Contains(":name test-function", f1)
	as.String("this is a test", f1.Documentation())

	c := a.NewContext()
	as.String("hello", f1.Apply(c, a.EmptyList))

	f4 := a.NewFunction(nil).WithMetadata(a.Properties{
		a.MetaType: f(99),
	}).(a.Function)

	as.String("function", f4.Type())
}

func TestMacro(t *testing.T) {
	as := assert.New(t)

	foo := a.NewKeyword("foo")

	m1 := a.NewFunction(nil).WithMetadata(a.Properties{
		a.MetaMacro: a.True,
		a.MetaName:  a.Name("orig"),
	}).(a.Function)

	ok, _ := a.IsMacro(m1)
	as.True(ok)

	ok, _ = a.IsMacro(a.NewFunction(nil))
	as.False(ok)

	k1 := a.NewKeyword("some_keyword")

	ok, _ = a.IsMacro(k1)
	as.False(ok)

	ok = a.IsSpecialForm(k1)
	as.False(ok)

	m2 := m1.WithMetadata(a.Properties{
		foo:        s("bar"),
		a.MetaName: a.Name("changed"),
	}).(a.Function)

	v, _ := m1.Metadata().Get(a.MetaMacro)
	as.True(v)
	v, _ = m2.Metadata().Get(a.MetaMacro)
	as.True(v)

	as.Contains(":type function", m1)

	v, _ = m1.Metadata().Get(a.MetaName)
	as.String("orig", v)

	v, _ = m2.Metadata().Get(a.MetaName)
	as.String("changed", v)

	v, _ = m2.Metadata().Get(foo)
	as.String("bar", v)

	v, _ = m1.Metadata().Get(foo)
	as.Nil(v)
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
