package api

const (
	// BadArity is thrown when a Function has a fixed arity
	BadArity = "expected %d argument(s), got %d"

	// BadMinimumArity is thrown when a Function has a minimum arity
	BadMinimumArity = "expected at least %d argument(s), got %d"

	// BadArityRange is thrown when a Function has an arity range
	BadArityRange = "expected between %d and %d arguments, got %d"
)

// ArityChecker is a ReflectedFunction that validates the arity of arguments
type ArityChecker func(Sequence) (int, bool)

func countUpTo(args Sequence, c int) int {
	if cnt, ok := args.(Counted); ok {
		return cnt.Count()
	}
	r := 0
	for s := args; r < c && s.IsSequence(); s = s.Rest() {
		r++
	}
	return r
}

// MakeArityChecker creates a fixed arity checker
func MakeArityChecker(arity int) ArityChecker {
	plusOne := arity + 1
	return func(args Sequence) (int, bool) {
		c := countUpTo(args, plusOne)
		return c, c == arity
	}
}

// AssertArity explodes if the arg count doesn't match provided arity
func AssertArity(args Sequence, arity int) int {
	c, ok := MakeArityChecker(arity)(args)
	if !ok {
		panic(ErrStr(BadArity, arity, c))
	}
	return c
}

// MakeMinimumArityChecker creates a minimum arity checker
func MakeMinimumArityChecker(arity int) ArityChecker {
	return func(args Sequence) (int, bool) {
		c := countUpTo(args, arity)
		return c, c >= arity
	}
}

// AssertMinimumArity explodes if the arg count isn't at least arity
func AssertMinimumArity(args Sequence, arity int) int {
	c, ok := MakeMinimumArityChecker(arity)(args)
	if !ok {
		panic(ErrStr(BadMinimumArity, arity, c))
	}
	return c
}

// MakeArityRangeChecker creates a ranged arity checker
func MakeArityRangeChecker(min int, max int) ArityChecker {
	maxPlusOne := max + 1
	return func(args Sequence) (int, bool) {
		c := countUpTo(args, maxPlusOne)
		return c, c >= min && c <= max
	}
}

// AssertArityRange explodes if the arg count isn't in the arity range
func AssertArityRange(args Sequence, min int, max int) int {
	c, ok := MakeArityRangeChecker(min, max)(args)
	if !ok {
		panic(ErrStr(BadArityRange, min, max, c))
	}
	return c
}
