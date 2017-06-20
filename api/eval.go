package api

// ExpectedApplicable is thrown when a Value is not Applicable
const ExpectedApplicable = "value does not support application: %s"

// Applicable is the standard signature for any Value that can be applied
// to a sequence of arguments
type Applicable interface {
	Apply(Context, Sequence) Value
}

// AssertApplicable will cast a Value into an Applicable or explode violently
func AssertApplicable(v Value) Applicable {
	if r, ok := v.(Applicable); ok {
		return r
	}
	panic(Err(ExpectedApplicable, v))
}
