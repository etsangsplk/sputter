package sputter_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	r "github.com/kode4food/sputter/reader"
	"github.com/stretchr/testify/assert"
)

func testCode(a *assert.Assertions, src string, expect s.Value) {
	ctx := s.NewContext()
	l := r.NewLexer(src)
	c := r.NewCoder(b.BuiltIns, l)
	a.Equal(expect, r.EvaluateCoder(ctx, c), src)
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
}

func TestBadArity(t *testing.T) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			err := "expected 1 argument(s), got 0"
			a.Equal(rec, err, "bad arity")
			return
		}
		a.Fail("bad arity didn't panic")
	}()

	testCode(a, `
		(defun identity [value] value)
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
