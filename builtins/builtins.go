package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

const (
	badArity      = "expected %d argument(s), got %d"
	badMinArity   = "expected at least %d argument(s), got %d"
	badArityRange = "expected between %d and %d arguments, got %d"
)

// Context is a special Context of built-in identifiers
var Context = a.NewContext()

func countIterator(i a.Iterator) int {
	var c = 0
	for _, ok := i.Next(); ok; _, ok = i.Next() {
		c++
	}
	return c
}

func argCount(args a.Iterable) int {
	if c, ok := args.(a.Countable); ok {
		return c.Count()
	}
	return countIterator(args.Iterate())
}

// AssertArity explodes if the arg count doesn't match provided arity
func AssertArity(args a.Iterable, arity int) {
	c := argCount(args)
	if c != arity {
		panic(fmt.Sprintf(badArity, arity, c))
	}
}

// AssertMinimumArity explodes if the arg count isn't at least arity
func AssertMinimumArity(args a.Iterable, arity int) {
	c := argCount(args)
	if c < arity {
		panic(fmt.Sprintf(badMinArity, arity, c))
	}
}

// AssertArityRange explodes if the arg count isn't in the arity range
func AssertArityRange(args a.Iterable, min int, max int) {
	c := argCount(args)
	if c < min || c > max {
		panic(fmt.Sprintf(badArityRange, min, max, c))
	}
}

func init() {
	Context.Put("T", a.True)
	Context.Put("nil", a.Nil)
	Context.Put("true", a.True)
	Context.Put("false", a.False)
}
