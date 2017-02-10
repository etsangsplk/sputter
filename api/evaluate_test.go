package api_test

import (
	"fmt"
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func evaluateToString(c *s.Context, v s.Value) string {
	result := s.Evaluate(c, v)
	if str, ok := result.(fmt.Stringer); ok {
		return str.String()
	}
	return result.(string)
}

func TestEvaluate(t *testing.T) {
	a := assert.New(t)

	hello := &s.Function{
		Name: "hello",
		Exec: func(c *s.Context, args s.Iterable) s.Value {
			iter := args.Iterate()
			arg, _ := iter.Next()
			value := evaluateToString(c, arg)
			return "Hello, " + value + "!"
		},
	}

	list := s.NewList(hello).Conj("World")
	result := s.Evaluate(s.NewContext(), list)

	a.Equal("Hello, World!", result.(string), "good hello")
}
