package api

import "fmt"

const (
	// BadArity is thrown when a Function has a fixed arity
	BadArity = "expected %d argument(s), got %d"

	// BadMinimumArity is thrown when a Function has a minimum arity
	BadMinimumArity = "expected at least %d argument(s), got %d"

	// BadArityRange is thrown when a Function has an arity range
	BadArityRange = "expected between %d and %d arguments, got %d"

	// ExpectedFunction is thrown when a Value is not a Function
	ExpectedFunction = "value is not a function"
)

// SequenceProcessor is the standard signature for a function that is
// capable of transforming or validating a Sequence
type SequenceProcessor func(Context, Sequence) Value

// Function is a Value that can be invoked
type Function struct {
	Name    Name
	Apply   SequenceProcessor
	Prepare SequenceProcessor
	Data    bool
}

// ResolveAsFunction either returns a Function as-is or tries
// to perform a lookup if the Value is a Symbol
func ResolveAsFunction(c Context, v Value) (*Function, bool) {
	if f, ok := v.(*Function); ok {
		return f, true
	}

	if s, ok := v.(*Symbol); ok {
		if sv, ok := s.Resolve(c); ok {
			if f, ok := sv.(*Function); ok {
				return f, true
			}
		}
	}

	return nil, false
}

func (f *Function) String() string {
	if f.Name != "" {
		return "(fn :name " + string(f.Name) + ")"
	}
	return fmt.Sprintf("(fn :addr %p)", &f)
}

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

// AssertFunction will cast a Value into a Function or explode violently
func AssertFunction(v Value) *Function {
	if r, ok := v.(*Function); ok {
		return r
	}
	panic(ExpectedFunction)
}
