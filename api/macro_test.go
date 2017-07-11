package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestMacro(t *testing.T) {
	as := assert.New(t)

	foo := a.NewKeyword("foo")

	m1 := a.NewMacro(nil).WithMetadata(a.Properties{
		a.MetaName: a.Name("orig"),
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

	as.Contains(":type macro", m1)

	v, _ = m1.Metadata().Get(a.MetaName)
	as.String("orig", v)

	v, _ = m2.Metadata().Get(a.MetaName)
	as.String("changed", v)

	v, _ = m2.Metadata().Get(foo)
	as.String("bar", v)

	v, _ = m1.Metadata().Get(foo)
	as.Nil(v)
}
