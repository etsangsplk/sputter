package sputter_test

import (
	"math/big"
	"testing"

	"fmt"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	i "github.com/kode4food/sputter/interpreter"
	"github.com/stretchr/testify/assert"
)

func testCodeWithContext(a *assert.Assertions, code string,
	expect s.Value, context *s.Context) {
	l := i.NewLexer(code)
	c := i.NewCoder(b.BuiltIns, l)
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
	`, s.EmptyList)
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
