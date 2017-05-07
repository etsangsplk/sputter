package native_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	r "github.com/kode4food/sputter/native"
)

func TestStructReflect(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()

	r1, ok := n1.Get(a.Name("s"))
	as.True(ok)
	as.String("hello", r1)

	r2, ok := n1.Get(a.Name("b"))
	as.True(ok)
	as.True(r2)

	r3, ok := n1.Get(a.Name("i"))
	as.True(ok)
	as.Number(42, r3)

	r4, ok := n1.Get(a.Name("f"))
	as.True(ok)
	as.Number(99.5, r4)

	r5, ok := n1.Get(a.Name("is-camel-case"))
	as.True(ok)
	as.String("I was camelCase", r5)

	r6, ok := n1.Get(a.Name("not-exported"))
	as.False(ok)
	as.Nil(r6)
}

func TestNestedStructReflect(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()

	r6, ok := n1.Get(a.Name("n"))
	as.True(ok)
	n2, ok := r6.(r.Value)
	as.True(ok)
	as.NotNil(n2)
	r7, ok := n2.Get(a.Name("nested"))
	as.True(ok)
	as.True(r7)

	r8, ok := n1.Get(a.Name("r"))
	as.True(ok)
	n3, ok := r8.(r.Value)
	as.True(ok)
	as.NotNil(n3)
	r9, ok := n3.Get(a.Name("s"))
	as.True(ok)
	as.String("nested", r9)
}

func TestBadStructReflect(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()

	r1, ok := n1.Get(a.Name("s"))
	as.True(ok)
	as.String("hello", r1)

	defer as.ExpectError(a.Err(r.BadConversionType, "[]int"))
	n1.Get(a.Name("a"))
}
