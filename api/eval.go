package api

// Eval evaluates a Value against a Context
func Eval(c *Context, v Value) Value {
	if e, ok := v.(Evaluable); ok {
		return e.Eval(c)
	}
	return v
}

// EvalSequence evaluates each element of the provided Sequence
func EvalSequence(c *Context, s Sequence) Value {
	var r Value = Nil
	i := s.Iterate()
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		r = Eval(c, v)
	}
	return r
}
