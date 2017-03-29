package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestNames(t *testing.T) {
	as := assert.New(t)

	n := a.Name("hello")
	as.Equal(a.Name("hello"), n.Name(), "Name Name() works")
}
