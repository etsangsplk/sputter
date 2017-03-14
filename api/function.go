package api

import "fmt"

const (
	// BadArity is thrown when a Function has a fixed arity
	BadArity = "expected %d argument(s), got %d"

	// BadMinimumArity is thrown when a Function has a minimum arity
	BadMinimumArity = "expected at least %d argument(s), got %d"

	// BadArityRange is thrown when a Function has an arity range
	BadArityRange = "expected between %d and %d arguments, got %d"

	// ExpectedApplicable is thrown when a Value is not Applicable
	ExpectedApplicable = "value does not support application"
)

// SequenceProcessor is the standard signature for a function that is
// capable of transforming or validating a Sequence
type SequenceProcessor func(Context, Sequence) Value

// Applicable is the standard signature for any Value that can have
// arguments applied to it
type Applicable interface {
	Apply(Context, Sequence) Value
}

// Function is a Value that can be invoked
type Function struct {
	Name Name
	Exec SequenceProcessor
}

// Macro is a Value that can be used to transform a Value
type Macro struct {
	Name Name
	Prep SequenceProcessor
	Data bool
}

// ResolveAsApplicable either returns an Applicable as-is or tries
// to perform a lookup if the Value is a Symbol
func ResolveAsApplicable(c Context, v Value) (Applicable, bool) {
	if f, ok := v.(Applicable); ok {
		return f, true
	}

	if s, ok := v.(*Symbol); ok {
		if sv, ok := s.Resolve(c); ok {
			if f, ok := sv.(Applicable); ok {
				return f, true
			}
		}
	}

	return nil, false
}

// Apply makes Function Applicable
func (f *Function) Apply(c Context, args Sequence) Value {
	return f.Exec(c, args)
}

func (f *Function) String() string {
	if f.Name != "" {
		return "(fn :name " + string(f.Name) + ")"
	}
	return fmt.Sprintf("(fn :addr %p)", &f)
}

// Apply makes Macro Applicable
func (m *Macro) Apply(c Context, args Sequence) Value {
	return m.Prep(c, args)
}

func (m *Macro) String() string {
	if m.Name != "" {
		return "(macro :name " + string(m.Name) + ")"
	}
	return fmt.Sprintf("(macro :addr %p)", &m)
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

// AssertApplicable will cast a Value into an Applicable or explode violently
func AssertApplicable(v Value) Applicable {
	if r, ok := v.(Applicable); ok {
		return r
	}
	panic(ExpectedApplicable)
}
