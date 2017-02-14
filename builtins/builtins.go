package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

const (
	incorrectArity      = "expected %d argument(s), got %d"
	incorrectMinArity   = "expected at least %d argument(s), got %d"
	incorrectArityRange = "expected between %d and %d arguments, got %d"
)

// BuiltIns is a special Context of built-in identifiers
var BuiltIns = a.NewContext()

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
	count := argCount(args)
	if count != arity {
		panic(fmt.Sprintf(incorrectArity, arity, count))
	}
}

// AssertMinimumArity explodes if the arg count isn't at least arity
func AssertMinimumArity(args a.Iterable, arity int) {
	count := argCount(args)
	if count < arity {
		panic(fmt.Sprintf(incorrectMinArity, arity, count))
	}
}

// AssertArityRange explodes if the arg count isn't in the arity range
func AssertArityRange(args a.Iterable, min int, max int) {
	count := argCount(args)
	if count < min || count > max {
		panic(fmt.Sprintf(incorrectArityRange, min, max, count))
	}
}

func init() {
	BuiltIns.Put("T", a.True)
	BuiltIns.Put("nil", a.Nil)
	BuiltIns.Put("true", a.True)
	BuiltIns.Put("false", a.False)
}
