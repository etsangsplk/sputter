package native_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestSingularMethodInvoke(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()
	v1, ok := n1.Get(a.Name("single"))
	as.True(ok)
	as.NotNil(v1)

	f1, ok := v1.(*a.Function)
	as.True(ok)
	as.NotNil(f1)

	c := a.NewContext()
	r1 := f1.Apply(c, &a.Vector{a.Str("hello")})
	as.NotNil(r1)
	as.Number(5, r1)
}

func TestVoidMethodInvoke(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()
	v1, ok := n1.Get(a.Name("void"))
	as.True(ok)
	as.NotNil(v1)

	f1, ok := v1.(*a.Function)
	as.True(ok)
	as.NotNil(f1)

	c := a.NewContext()
	r1 := f1.Apply(c, &a.Vector{a.Str("clobbered"), a.NewFloat(-99)})
	as.Nil(r1)

	r2, ok := n1.Get(a.Name("s"))
	as.True(ok)
	as.String("clobbered", r2)

	r3, ok := n1.Get(a.Name("f"))
	as.True(ok)
	as.Number(-99, r3)
}

func TestPluralMethodInvoke(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()
	v1, ok := n1.Get(a.Name("plural"))
	as.True(ok)
	as.NotNil(v1)

	f1, ok := v1.(*a.Function)
	as.True(ok)
	as.NotNil(f1)

	c := a.NewContext()
	r1 := f1.Apply(c, &a.Vector{a.Str("one"), a.NewFloat(4), a.Str("eight")})
	as.NotNil(r1)

	r2 := r1.(a.Vector)
	as.Number(3, r2[0])
	as.Number(4, r2[1])
	as.Number(5, r2[2])
}
