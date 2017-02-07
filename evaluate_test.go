package main

import (
	"math/big"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func testCodeWithContext(a *assert.Assertions, code string,
	expect Value, context *Context) {
	l := NewLexer(code)
	c := NewCoder(l)
	a.Equal(expect, EvaluateCoder(context, c), code)
}

func testCode(a *assert.Assertions, code string, expect Value) {
	c := Builtins.Child()
	testCodeWithContext(a, code, expect, c)
}

func evaluateToString(c *Context, v Value) string {
	result := Evaluate(c, v)
	if str, ok := result.(fmt.Stringer); ok {
		return str.String()
	}
	return result.(string)
}

func TestBasics(t *testing.T) {
	a := assert.New(t)
	testCode(a, "(+ 1 1)", big.NewFloat(2.0))
	testCode(a, "(* 4 4)", big.NewFloat(16.0))
	testCode(a, "(+ 5 4)", big.NewFloat(9.0))
	testCode(a, "(* 12 3)", big.NewFloat(36.0))
	testCode(a, "(- 10 4)", big.NewFloat(6.0))
	testCode(a, "(- 10 4 2)", big.NewFloat(4.0))
	testCode(a, "(/ 10 2)", big.NewFloat(5.0))
	testCode(a, "(/ 10 2 5)", big.NewFloat(1.0))
}

func TestNested(t *testing.T) {
	a := assert.New(t)
	testCode(a, "(/ 10 (- 5 3))", big.NewFloat(5.0))
	testCode(a, "(* 5 (- 5 3))", big.NewFloat(10.0))
	testCode(a, "(/ 10 (/ 6 3))", big.NewFloat(5.0))
}

func TestEvaluable(t *testing.T) {
	a := assert.New(t)
	c := Builtins.Child()

	hello := &Function{"hello", func(c *Context, args Iterable) Value {
		iter := args.Iterate()
		arg, _ := iter.Next()
		value := evaluateToString(c, arg)
		return "Hello, " + value + "!"
	}}

	c.Put("hello", hello)
	c.Put("name", "Bob")

	testCodeWithContext(a, `(hello "World")`, "Hello, World!", c)
	testCodeWithContext(a, `(hello name)`, "Hello, Bob!", c)
}

func TestEvaluate(t *testing.T) {
	a := assert.New(t)

	hello := &Function{"hello", func(c *Context, args Iterable) Value {
		iter := args.Iterate()
		arg, _ := iter.Next()
		value := evaluateToString(c, arg)
		return "Hello, " + value + "!"
	}}

	list := NewList(hello).Conj("World")
	result := Evaluate(Builtins, list)

	a.Equal("Hello, World!", result.(string), "good hello")
}

func TestFunction(t *testing.T) {
	a := assert.New(t)

	testCode(a, `
		(defun say-hello [] "Hello, World!")
		(say-hello)
	`, "Hello, World!")

	testCode(a, `
		(defun identity [value] value)
		(identity "foo")
	`, "foo")

	testCode(a, `
	  (defun identity [value] value)
		(print '(identity "hello"))
		(identity)
	`, EmptyList)
}

func TestVariables(t *testing.T) {
	a := assert.New(t)

	testCode(a, `
		(defvar foo "bar")
		foo
	`, "bar")

	testCode(a, `
		(defun return-local []
			(let [foo "local"] foo))
		(return-local)
	`, "local")
}
