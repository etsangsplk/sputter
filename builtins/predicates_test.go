package builtins_test

import (
	"fmt"
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestPredicates(t *testing.T) {
	testCode(t, `(eq true true)`, a.True)
	testCode(t, `(eq true false)`, a.False)
	testCode(t, `(eq false false)`, a.True)
	testCode(t, `(eq 1 1)`, a.False)

	testCode(t, `(!eq true true)`, a.False)
	testCode(t, `(!eq true false)`, a.True)
	testCode(t, `(!eq false false)`, a.False)
	testCode(t, `(!eq 1 1)`, a.True)

	testCode(t, `(nil? ())`, a.True)
	testCode(t, `(nil? () () ())`, a.True)
	testCode(t, `(nil? () nil)`, a.True)
	testCode(t, `(nil? false)`, a.False)
	testCode(t, `(nil? false () nil)`, a.False)
	testCode(t, `(nil? nil)`, a.True)

	testCode(t, `(nil? "hello")`, a.False)
	testCode(t, `(nil? '(1 2 3))`, a.False)
	testCode(t, `(nil? () nil "hello")`, a.False)

	testBadCode(t, `(nil?)`, fmt.Sprintf(a.BadMinimumArity, 1, 0))
}
