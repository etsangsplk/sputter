package builtins_test

import (
	"fmt"
	"testing"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {
	testCode(t, `
		(defun say-hello [] "Hello, World!")
		(say-hello)
	`, "Hello, World!")

	testCode(t, `
		(defun identity [value] value)
		(identity "foo")
	`, "foo")
}

func TestBadArity(t *testing.T) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Sprintf(b.BadArity, 1, 0)
			a.Equal(err, rec, "bad arity")
			return
		}
		a.Fail("bad arity didn't panic")
	}()

	testCode(t, `
		(defun identity [value] value)
		(identity)
	`, s.Nil)
}
