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
	as.Identical(w, n.WrappedValue())

	v, ok := n.Metadata().Get(tk)
	as.True(ok)
	as.True(v)
}
