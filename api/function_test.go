package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {
	a := assert.New(t)

	f := &s.Function{
		Name: "test-function",
		Apply: func(c s.Context, args s.Sequence) s.Value {
			return "hello"
		},
	}

	a.Equal("test-function", f.String(), "string returned")

	c := s.NewContext()
	a.Equal("hello", f.Apply(c, s.Nil), "function executes")
}
