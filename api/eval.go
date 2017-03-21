package api

// Evaluable can be evaluated against a Context
type Evaluable interface {
	Eval(c Context) Value
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
	var r Value = Nil
	for i := s; i.IsSequence(); i = i.Rest() {
		r = Eval(c, i.First())
	}
	return r
}
