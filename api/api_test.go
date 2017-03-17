package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestTruthy(t *testing.T) {
	as := assert.New(t)

	as.True(a.Truthy(a.True), "API True is Truthy")
	as.True(a.Truthy(true), "true is Truthy")
	as.True(a.Truthy(a.NewList("Hello")), "Non-Empty List Is Truthy")
	as.True(a.Truthy("hello"), "String is Truthy")

	as.False(a.Truthy(a.Nil), "API Nil is not Truthy")
	as.False(a.Truthy(nil), "nil is not Truthy")
	as.False(a.Truthy(a.False), "API False is not Truthy")
	as.False(a.Truthy(false), "false is not Truthy")
}
