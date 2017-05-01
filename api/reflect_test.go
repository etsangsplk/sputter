package api_test

import (
	"bytes"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

type reflectNestedStruct struct {
	Nested bool
}

type reflectStruct struct {
	S string
	B bool
	I int32
	F float32
	A []int
	a bool
	N reflectNestedStruct
	R *reflectStruct
}

func getTestReflectStruct() a.Native {
	return a.NewNative(&reflectStruct{
		S: "hello",
		B: true,
		I: 42,
		F: 99.5,
		A: []int{1, 2, 3},
		a: true,
		N: reflectNestedStruct{
			Nested: true,
		},
		R: &reflectStruct{
			S: "nested",
		},
	})
}

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

func TestStructReflect(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()

	r1, ok := n1.Get(a.Name("S"))
	as.True(ok)
	as.String("hello", r1)

	r2, ok := n1.Get(a.Name("B"))
	as.True(ok)
	as.True(r2)

	r3, ok := n1.Get(a.Name("I"))
	as.True(ok)
	as.Number(42, r3)

	r4, ok := n1.Get(a.Name("F"))
	as.True(ok)
	as.Number(99.5, r4)

	r5, ok := n1.Get(a.Name("a"))
	as.False(ok)
	as.Nil(r5)
}

func TestNestedStructReflect(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()

	r6, ok := n1.Get(a.Name("N"))
	as.True(ok)
	n2, ok := r6.(a.Native)
	as.True(ok)
	as.NotNil(n2)
	r7, ok := n2.Get(a.Name("Nested"))
	as.True(ok)
	as.True(r7)

	r8, ok := n1.Get(a.Name("R"))
	as.True(ok)
	n3, ok := r8.(a.Native)
	as.True(ok)
	as.NotNil(n3)
	r9, ok := n3.Get(a.Name("S"))
	as.True(ok)
	as.String("nested", r9)
}

func TestBadStructReflect(t *testing.T) {
	as := assert.New(t)

	n1 := getTestReflectStruct()

	r1, ok := n1.Get(a.Name("S"))
	as.True(ok)
	as.String("hello", r1)

	defer expectError(as, a.Err(a.BadConversionType, "[]int"))
	n1.Get(a.Name("A"))
}
