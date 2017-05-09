package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestMacro(t *testing.T) {
	as := assert.New(t)

	foo := a.NewKeyword("foo")

	m1 := a.NewMacro(nil).WithMetadata(a.Metadata{
		a.MetaName: a.Name("orig"),
	}).(a.Macro)

	m2 := m1.WithMetadata(a.Metadata{
		foo:        s("bar"),
		a.MetaName: a.Name("changed"),
	}).(a.Macro)

	as.True(m1.Metadata()[a.MetaMacro])
	as.True(m2.Metadata()[a.MetaMacro])
	as.True(m1.DataMode())
	as.True(m2.DataMode())

	as.Contains(":type macro", m1)
	as.String("orig", m1.Metadata()[a.MetaName])
	as.String("changed", m2.Metadata()[a.MetaName])
	as.String("bar", m2.Metadata()[foo])
	as.Nil(m1.Metadata()[foo])
}
