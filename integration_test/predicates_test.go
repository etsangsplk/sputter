package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	_ "github.com/kode4food/sputter/core"
)

func TestPredicates(t *testing.T) {
	testCode(t, `(eq true true true)`, a.True)
	testCode(t, `(eq true false true)`, a.False)
	testCode(t, `(eq false false false)`, a.True)
	testCode(t, `(eq 1 1)`, a.False)

	testCode(t, `(!eq true true true)`, a.False)
	testCode(t, `(!eq true false)`, a.True)
	testCode(t, `(!eq false false)`, a.False)
	testCode(t, `(!eq 1 1)`, a.True)

	testCode(t, `(nil? nil)`, a.True)
	testCode(t, `(nil? nil nil nil)`, a.True)
	testCode(t, `(nil? () nil)`, a.False)
	testCode(t, `(nil? false)`, a.False)
	testCode(t, `(nil? false () nil)`, a.False)

	testCode(t, `(nil? "hello")`, a.False)
	testCode(t, `(nil? '(1 2 3))`, a.False)
	testCode(t, `(nil? () nil "hello")`, a.False)

	testCode(t, `(keyword? :hello)`, a.True)
	testCode(t, `(!keyword? :hello)`, a.False)
	testCode(t, `(keyword? 99)`, a.False)
	testCode(t, `(!keyword? 99)`, a.True)

	testBadCode(t, `(nil?)`, a.ErrStr(b.ExpectedArguments, "[first & rest]"))
}
