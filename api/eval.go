package api

// ExpectedApplicable is thrown when a Value is not Applicable
const ExpectedApplicable = "value does not support application: %s"

// Evaluable can be evaluated against a Context
type Evaluable interface {
	Eval(Context) Value
}

// Applicable is the standard signature for any Value that can have
// arguments applied to it
type Applicable interface {
	Value
	Apply(Context, Sequence) Value
}

// Eval evaluates a Value against a Context
func Eval(c Context, v Value) Value {
	if e, ok := v.(Evaluable); ok {
		return e.Eval(c)
	}
	return v
}

// EvalSequence evaluates each element of the provided Sequence
func EvalSequence(c Context, s Sequence) Value {
	r := Nil
	for i := s; i.IsSequence(); i = i.Rest() {
		r = Eval(c, i.First())
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
