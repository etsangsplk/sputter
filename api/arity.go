package api

const (
	// BadArity is thrown when a Function has a fixed arity
	BadArity = "expected %d argument(s), got %d"

	// BadMinimumArity is thrown when a Function has a minimum arity
	BadMinimumArity = "expected at least %d argument(s), got %d"

	// BadArityRange is thrown when a Function has an arity range
	BadArityRange = "expected between %d and %d arguments, got %d"
)

// ArityChecker is a function that validates the arity of arguments
type ArityChecker func(Values) (int, bool)

// AssertArity explodes if the arg count doesn't match provided arity
func AssertArity(args Values, arity int) int {
	c := len(args)
	if c != arity {
		panic(ErrStr(BadArity, arity, c))
	}
	return c
}

// AssertMinimumArity explodes if the arg count isn't at least arity
func AssertMinimumArity(args Values, arity int) int {
	c := len(args)
	if c < arity {
		panic(ErrStr(BadMinimumArity, arity, c))
	}
	return c
}

// AssertArityRange explodes if the arg count isn't in the arity range
func AssertArityRange(args Values, min int, max int) int {
	c := len(args)
	if c < min || c > max {
		panic(ErrStr(BadArityRange, min, max, c))
	}
	return c
}
