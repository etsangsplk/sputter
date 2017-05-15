package native_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	n "github.com/kode4food/sputter/native"
)

func someTestFunc(arg1, arg2 float64, b bool) (bool, float64) {
	return b, arg1 + arg2
}

func withValuesFunc(arg1 a.Str, arg2 a.Bool) a.Sequence {
	return a.NewVector(arg1, arg2)
}

func TestReflectFunc(t *testing.T) {
	as := assert.New(t)

	n1 := n.New(someTestFunc).(a.Applicable)
	c := a.NewContext()
	r1 := n1.Apply(c, a.NewVector(a.NewFloat(20), a.NewFloat(30), a.True))

	v1, ok := r1.(a.Vector)
	as.True(ok)
	as.NotNil(v1)

	v2, ok := v1.ElementAt(0)
	as.True(ok)
	as.True(v2)

	v3, ok := v1.ElementAt(1)
	as.True(ok)
	as.Number(50, v3)
}

func TestValuesFunc(t *testing.T) {
	as := assert.New(t)

	n1 := n.New(withValuesFunc).(a.Applicable)
	c := a.NewContext()
	r1 := n1.Apply(c, a.NewVector(a.Str("hello"), a.True))

	v1, ok := r1.(a.Vector)
	as.True(ok)
	as.NotNil(v1)

	v2, ok := v1.ElementAt(0)
	as.True(ok)
	as.String("hello", v2)

	v3, ok := v1.ElementAt(1)
	as.True(ok)
	as.True(v3)
}
