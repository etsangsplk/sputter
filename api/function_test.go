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
		Exec: func(c *s.Context, args s.Sequence) s.Value {
			return "hello"
		},
	}

	a.Equal("test-function", f.String(), "string returned")
	a.Equal("hello", f.Exec(s.NewContext(), s.Nil), "function executes")
}
