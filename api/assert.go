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

	// ExpectedList is thrown when a Value is not a Cons cell
	ExpectedList = "value is not a list"

	// ExpectedSequence is thrown when a Value is not a Sequence
	ExpectedSequence = "value is not a list or vector"

	// ExpectedUnqualified is thrown when a Value is not a unqualified Symbol
	ExpectedUnqualified = "value is not a symbol"

	// ExpectedNumeric is thrown when a Value is not a Number
	ExpectedNumeric = "value is not numeric"

	// ExpectedFunction is thrown when a Value is not a Function
	ExpectedFunction = "value is not a function"

	// ExpectedName is thrown when a Value is not a Name
	ExpectedName = "value is not a name"
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

// AssertCons will cast a Value into a Cons or explode violently
func AssertCons(v Value) *Cons {
	if r, ok := v.(*Cons); ok {
		return r
	}
	panic(ExpectedCons)
}

// AssertSequence will cast a Value into a Sequence or explode violently
func AssertSequence(v Value) Sequence {
	if r, ok := v.(Sequence); ok {
		return r
	}
	panic(ExpectedSequence)
}

// AssertUnqualified will cast a Value into a Symbol and explode
// violently if it's qualified with a domain
func AssertUnqualified(v Value) *Symbol {
	if r, ok := v.(*Symbol); ok {
		if r.Domain == LocalDomain {
			return r
		}
	}
	panic(ExpectedUnqualified)
}

// AssertNumeric will cast a Value into a Numeric or explode violently
func AssertNumeric(v Value) *big.Float {
	if r, ok := v.(*big.Float); ok {
		return r
	}
	panic(ExpectedNumeric)
}

// AssertFunction will cast a Value into a Function or explode violently
func AssertFunction(v Value) *Function {
	if r, ok := v.(*Function); ok {
		return r
	}
	panic(ExpectedFunction)
}

// AssertName will cast a Value to a Name or explode violently
func AssertName(v Value) Name {
	if r, ok := v.(Name); ok {
		return r
	}
	panic(ExpectedName)
}
