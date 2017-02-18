package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

const (
	// BadArity is thrown when a Function has a fixed arity
	BadArity = "expected %d argument(s), got %d"

	// BadMinArity is thrown when a Function has a minimum arity
	BadMinArity = "expected at least %d argument(s), got %d"

	// BadArityRange is thrown when a Function has an arity range
	BadArityRange = "expected between %d and %d arguments, got %d"
)

// Context is a special Context of built-in identifiers
var Context = a.NewContext()

// AssertArity explodes if the arg count doesn't match provided arity
func AssertArity(args a.Sequence, arity int) {
	c := a.Count(args)
	if c != arity {
		panic(fmt.Sprintf(BadArity, arity, c))
	}
}

// AssertMinimumArity explodes if the arg count isn't at least arity
func AssertMinimumArity(args a.Sequence, arity int) {
	c := a.Count(args)
	if c < arity {
		panic(fmt.Sprintf(BadMinArity, arity, c))
	}
}

// AssertArityRange explodes if the arg count isn't in the arity range
func AssertArityRange(args a.Sequence, min int, max int) {
	c := a.Count(args)
	if c < min || c > max {
		panic(fmt.Sprintf(BadArityRange, min, max, c))
	}
}

func init() {
	Context.Put("nil", a.Nil)
	Context.Put("true", a.True)
	Context.Put("false", a.False)
}
