package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

const (
	incorrectArity = "expected %d argument(s), got %d"
	incorrectMinArity = "expected at least %d argument(s), got %d"
	incorrectArityRange = "expected between %d and %d arguments, got %d"
)

// BuiltIns is a special Context of built-in identifiers
var BuiltIns = a.NewContext()

func argCount(args a.Iterable) int {
	if c, ok := args.(a.Countable); ok {
		return c.Count()
	}

	var count = 0
	iter := args.Iterate()
	for _, ok := iter.Next(); ok; _, ok = iter.Next() {
		count++
	}
	return count
}

func AssertArity(args a.Iterable, arity int) {
	count := argCount(args)
	if count != arity {
		panic(fmt.Sprintf(incorrectArity, arity, count))
	}
}

func AssertMinimumArity(args a.Iterable, arity int) {
	count := argCount(args)
	if count < arity {
		panic(fmt.Sprintf(incorrectMinArity, arity, count))
	}
}

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
