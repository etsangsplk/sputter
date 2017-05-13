package api

// ExpectedApplicable is thrown when a Value is not Applicable
const ExpectedApplicable = "value does not support application: %s"

// Evaluable can be evaluated against a Context
type Evaluable interface {
	Eval(Context) Value
}

// Applicable is the standard signature for any Value that can be applied
// to a sequence of arguments
type Applicable interface {
	Apply(Context, Sequence) Value
}

// Eval evaluates a Value against a Context
func Eval(c Context, v Value) Value {
	if e, ok := v.(Evaluable); ok {
		return e.Eval(c)
	}
	return v
}

// Apply will apply an Applicable to the specified arguments
func Apply(c Context, a Applicable, args Sequence) Value {
	//if IsMacro(a) {
	return a.Apply(c, args)
	//}
	//if cnt, ok := args.(Counted); ok {
	//	l := cnt.Count()
	//	v := make(Vector, l)
	//	for i, e := 0, args; i < l; i++ {
	//		v[i] = Eval(c, e.First())
	//		e = e.Rest()
	//	}
	//	return a.Apply(c, v)
	//}
	//v := Vector{}
	//for e := args; e.IsSequence(); e = e.Rest() {
	//	v = append(v, Eval(c, e.First()))
	//}
	//return a.Apply(c, v)
}

// EvalSequence evaluates each element of the provided Sequence
func EvalSequence(c Context, s Sequence) Value {
	var r Value = Nil
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
