package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestTruthy(t *testing.T) {
	a := assert.New(t)

	a.True(s.Truthy(s.True), "API True is Truthy")
	a.True(s.Truthy(true), "true is Truthy")
	a.True(s.Truthy(s.NewList("Hello")), "Non-Empty List Is Truthy")
	a.True(s.Truthy("hello"), "String is Truthy")

	a.False(s.Truthy(s.Nil), "API Nil is not Truthy")
	a.False(s.Truthy(nil), "nil is not Truthy")
	a.False(s.Truthy(s.False), "API False is not Truthy")
	a.False(s.Truthy(false), "false is not Truthy")
}
