package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestMacro(t *testing.T) {
	as := assert.New(t)

	foo := a.NewKeyword("foo")

	m1 := a.NewMacro(nil).WithMetadata(a.Metadata{
		a.MetaName: "orig",
	}).(a.Macro)

	m2 := m1.WithMetadata(a.Metadata{
		foo:        "bar",
		a.MetaName: "changed",
	}).(a.Macro)

	as.True(m1.Metadata()[a.MetaMacro].(bool))
	as.True(m2.Metadata()[a.MetaMacro].(bool))
	as.True(m1.DataMode())
	as.True(m2.DataMode())

	as.Equal("orig", m1.Metadata()[a.MetaName])
	as.Equal("changed", m2.Metadata()[a.MetaName])
	as.Equal("bar", m2.Metadata()[foo])
	as.Nil(m1.Metadata()[foo])
}
