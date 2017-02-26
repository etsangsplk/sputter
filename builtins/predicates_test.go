package builtins_test

import (
	"fmt"
	"testing"

	s "github.com/kode4food/sputter/api"
)

func TestPredicates(t *testing.T) {
	testCode(t, `(eq true true)`, s.True)
	testCode(t, `(eq true false)`, s.False)
	testCode(t, `(eq false false)`, s.True)
	testCode(t, `(eq 1 1)`, s.False)

	testCode(t, `(nil? ())`, s.True)
	testCode(t, `(nil? () () ())`, s.True)
	testCode(t, `(nil? () nil)`, s.True)	
	testCode(t, `(nil? false)`, s.False)
	testCode(t, `(nil? false () nil)`, s.False)
	testCode(t, `(nil? nil)`, s.True)

	testCode(t, `(nil? "hello")`, s.False)
	testCode(t, `(nil? '(1 2 3))`, s.False)
	testCode(t, `(nil? () nil "hello")`, s.False)

	testBadCode(t, `(nil?)`, fmt.Sprintf(s.BadMinimumArity, 1, 0))
}
