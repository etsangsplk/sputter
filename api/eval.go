package api

// ExpectedApplicable is thrown when a Value is not Applicable
const ExpectedApplicable = "value does not support application: %s"

// MakeExpression marks Types that can be converted to Expression
type MakeExpression interface {
	Expression() Value
}

// Expression marks a Type as being a transformative Expression
type Expression interface {
	IsExpression() bool
}

// Applicable is the standard signature for any Value that can be applied
// to a sequence of arguments
type Applicable interface {
	Apply(Context, Sequence) Value
}

// EvalBlock evaluates each element of the provided Sequence
func EvalBlock(c Context, s Sequence) Value {
	var r Value = Nil
	for i := s; i.IsSequence(); i = i.Rest() {
		r = i.First().Eval(c)
	}
	return r
}

// AssertApplicable will cast a Value into an Applicable or explode violently
func AssertApplicable(v Value) Applicable {
	if r, ok := v.(Applicable); ok {
		return r
	}
	panic(Err(ExpectedApplicable, v))
}
