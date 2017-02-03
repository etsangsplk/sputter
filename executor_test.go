package sputter

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testCodeWithContext(a *assert.Assertions, code string,
	expect Value, context *Context) {
	l := NewLexer(code)
	c := NewCoder(l)
	e := NewExecutor(c)
	a.Equal(expect, e.Exec(context), code)
}

func testCode(a *assert.Assertions, code string, expect Value) {
	testCodeWithContext(a, code, expect, Builtins)
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
}

func TestResolvable(t *testing.T) {
	a := assert.New(t)
	c := Builtins.Child()

	hello := &Function{func(c *Context, args *List) Value {
		value := args.value.(string)
		return "Hello, " + value + "!"
	}}

	c.Put("hello", hello)
	testCodeWithContext(a, `(hello "World")`, "Hello, World!", c)
}

func TestEvaluate(t *testing.T) {
	a := assert.New(t)

	hello := &Function{func(c *Context, args *List) Value {
		value := args.value.(string)
		return "Hello, " + value + "!"
	}}

	list := NewList(hello).Conj("World")
	result := Evaluate(Builtins, list)

	a.Equal("Hello, World!", result.(string), "good hello")
}
