package api_test

import (
	"strings"
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

type anon struct{}

func TestString(t *testing.T) {
	as := assert.New(t)

	as.Equal("true", a.String(a.True), "true stringifies correctly")
	as.Equal("false", a.String(a.False), "true stringifies correctly")
	as.Equal("nil", a.String(a.Nil), "nil stringifies correctly")
	as.Equal("nil", a.String(nil), "nil stringifies correctly")
	as.True(strings.HasPrefix(a.String(&anon{}), "(<anon>"))
}
