package native_test

import (
	"bytes"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	r "github.com/kode4food/sputter/native"
)

type reflectNestedStruct struct {
	Nested bool
}

type reflectStruct struct {
	S           string
	B           bool
	I           int32
	F           float32
	A           []int
	N           reflectNestedStruct
	R           *reflectStruct
	notExported bool
	IsCamelCase string
}

func (r *reflectStruct) Void(s string, f float32) {
	r.S = s
	r.F = f
}

func (r *reflectStruct) Single(n string) int {
	return len(n)
}

func (r *reflectStruct) Plural(s1 string, i int, s2 string) (int, int, int) {
	return len(s1), i, len(s2)
}

func getTestReflectStruct() r.Value {
	return r.NewValue(&reflectStruct{
		S: "hello",
		B: true,
		I: 42,
		F: 99.5,
		A: []int{1, 2, 3},
		N: reflectNestedStruct{
			Nested: true,
		},
		R: &reflectStruct{
			S: "nested",
		},
		notExported: true,
		IsCamelCase: "I was camelCase",
	})
}

func TestReflect(t *testing.T) {
	as := assert.New(t)

	w := bytes.NewBufferString("")
	tk := a.NewKeyword("test")
	n := r.NewValue(w).WithMetadata(a.Metadata{
		tk: a.True,
	}).(r.Value)

	as.String("*bytes.buffer", n.(a.Typed).Type())
	as.Contains(":type *bytes.buffer", n)
	as.Identical(w, n.Wrapped())

	v, ok := n.Metadata().Get(tk)
	as.True(ok)
	as.True(v)
}
