package interpreter_test

import (
	"testing"

	"fmt"

	s "github.com/kode4food/sputter/api"
	i "github.com/kode4food/sputter/interpreter"
	"github.com/stretchr/testify/assert"
)

func testCodeWithContext(a *assert.Assertions, code string,
	expect s.Value, context *s.Context) {
	l := i.NewLexer(code)
	c := i.NewCoder(s.NewContext(), l)
	a.Equal(expect, i.EvaluateCoder(context, c), code)
}

func testCode(a *assert.Assertions, code string, expect s.Value) {
	c := s.NewContext()
	testCodeWithContext(a, code, expect, c)
}

func evaluateToString(c *s.Context, v s.Value) string {
	result := i.Evaluate(c, v)
	if str, ok := result.(fmt.Stringer); ok {
		return str.String()
	}
	return result.(string)
}

func TestEvaluable(t *testing.T) {
	a := assert.New(t)
	c := s.NewContext()

	hello := &s.Function{
		Name: "hello",
		Exec: func(c *s.Context, args s.Iterable) s.Value {
			iter := args.Iterate()
			arg, _ := iter.Next()
			value := evaluateToString(c, arg)
			return "Hello, " + value + "!"
		},
	}

	c.Put("hello", hello)
	c.Put("name", "Bob")

	testCodeWithContext(a, `(hello "World")`, "Hello, World!", c)
	testCodeWithContext(a, `(hello name)`, "Hello, Bob!", c)
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
	result := i.Evaluate(s.NewContext(), list)

	a.Equal("Hello, World!", result.(string), "good hello")
}
