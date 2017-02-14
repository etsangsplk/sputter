package api

// Evaluate a Value against a Context
func Evaluate(c *Context, v Value) Value {
	if e, ok := v.(Evaluable); ok {
		return e.Evaluate(c)
	}
	return v
}

// EvaluateIterator evaluates each element of the provided Iterator
func EvaluateIterator(c *Context, i Iterator) Value {
	var r Value = EmptyList
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		r = Evaluate(c, v)
	}
	return r
}
