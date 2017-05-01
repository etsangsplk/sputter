package api_test

import (
	"bytes"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestReflect(t *testing.T) {
	as := assert.New(t)

	w := bytes.NewBufferString("")
	tk := a.NewKeyword("test")
	n := a.NewNative(w).WithMetadata(a.Metadata{
		tk: a.True,
	}).(a.Native)

	as.String("*bytes.Buffer", n.(a.Typed).Type())
	as.Contains(":type *bytes.Buffer", n)
	as.Identical(w, n.NativeValue())

	v, ok := n.Metadata().Get(tk)
	as.True(ok)
	as.True(v)
}

type reflectStruct struct {
	S string
	B bool
	I int32
	F float32
	A []int
	a bool
}

func TestStructReflect(t *testing.T) {
	as := assert.New(t)

	s := &reflectStruct{
		S: "hello",
		B: true,
		I: 42,
		F: 99.5,
		A: []int{1, 2, 3},
		a: true,
	}

	n := a.NewNative(s)

	r1, ok := n.Get(a.Name("S"))
	as.True(ok)
	as.String("hello", r1)

	r2, ok := n.Get(a.Name("B"))
	as.True(ok)
	as.True(r2)

	r3, ok := n.Get(a.Name("I"))
	as.True(ok)
	as.Number(42, r3)

	r4, ok := n.Get(a.Name("F"))
	as.True(ok)
	as.Number(99.5, r4)

	r5, ok := n.Get(a.Name("a"))
	as.False(ok)
	as.Nil(r5)

	defer expectError(as, a.Err(a.BadConversionType, "[]int"))
	n.Get(a.Name("A"))
}
