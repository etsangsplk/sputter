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

	f1, ok := v1.(a.Function)
	as.True(ok)
	as.NotNil(f1)

	c := a.Variables{}
	r1 := f1.Apply(c, a.NewVector(a.Str("hello")))
	as.NotNil(r1)
	as.Number(5, r1)
}

func TestVoidMethodInvoke(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()
	v1, ok := n1.Get(a.Name("void"))
	as.True(ok)
	as.NotNil(v1)

	f1, ok := v1.(a.Function)
	as.True(ok)
	as.NotNil(f1)

	c := a.Variables{}
	r1 := f1.Apply(c, a.NewVector(a.Str("clobbered"), a.NewFloat(-99)))
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

	f1, ok := v1.(a.Function)
	as.True(ok)
	as.NotNil(f1)

	c := a.Variables{}
	r1 := f1.Apply(c, a.NewVector(a.Str("one"), a.NewFloat(4), a.Str("eight")))
	as.NotNil(r1)

	r2 := r1.(a.Vector)
	r3, _ := r2.ElementAt(0)
	r4, _ := r2.ElementAt(1)
	r5, _ := r2.ElementAt(2)

	as.Number(3, r3)
	as.Number(4, r4)
	as.Number(5, r5)
}
