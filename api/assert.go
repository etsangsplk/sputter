package api

import (
	"fmt"
	"math/big"
)

const (
	// BadArity is thrown when a Function has a fixed arity
	BadArity = "expected %d argument(s), got %d"

	// BadMinimumArity is thrown when a Function has a minimum arity
	BadMinimumArity = "expected at least %d argument(s), got %d"

	// BadArityRange is thrown when a Function has an arity range
	BadArityRange = "expected between %d and %d arguments, got %d"

	// ExpectedCons is thrown when a Value is not a Cons cell
	ExpectedCons = "value is not a cons cell"

	// ExpectedSequence is thrown when a Value is not a Sequence
	ExpectedSequence = "value is not a list or vector"

	// ExpectedSymbol is thrown when a Value is not a Symbol
	ExpectedSymbol = "value is not a symbol"

	// ExpectedNumeric is thrown when a Value is not a Number
	ExpectedNumeric = "value is not numeric"

	// ExpectedFunction is thrown when a Value is not a Function
	ExpectedFunction = "value is not a function"
)

// AssertArity explodes if the arg count doesn't match provided arity
func AssertArity(args Sequence, arity int) {
	c := Count(args)
	if c != arity {
		panic(fmt.Sprintf(BadArity, arity, c))
	}
}

// AssertMinimumArity explodes if the arg count isn't at least arity
func AssertMinimumArity(args Sequence, arity int) {
	c := Count(args)
	if c < arity {
		panic(fmt.Sprintf(BadMinimumArity, arity, c))
	}
}

// AssertArityRange explodes if the arg count isn't in the arity range
func AssertArityRange(args Sequence, min int, max int) {
	c := Count(args)
	if c < min || c > max {
		panic(fmt.Sprintf(BadArityRange, min, max, c))
	}
}

// AssertSequence will cast a Value into a Sequence or explode violently
func AssertSequence(v Value) Sequence {
	if r, ok := v.(Sequence); ok {
		return r
	}
	panic(ExpectedSequence)
}

// AssertSymbol will cast a Value into a Symbol or explode violently
func AssertSymbol(v Value) *Symbol {
	if r, ok := v.(*Symbol); ok {
		return r
	}
	panic(ExpectedSymbol)
}

// AssertNumeric will cast a Value into a Numeric or explode violently
func AssertNumeric(v Value) *big.Float {
	if r, ok := v.(*big.Float); ok {
		return r
	}
	panic(ExpectedNumeric)
}
