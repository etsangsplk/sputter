package api

// Eval evaluates a Value against a Context
func Eval(c *Context, v Value) Value {
	if e, ok := v.(Evaluable); ok {
		return e.Eval(c)
	}
	return v
}

// EvalIterator evaluates each element of the provided Iterator
func EvalIterator(c *Context, i Iterator) Value {
	var r Value = Nil
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		r = Eval(c, v)
	}
	return r
}
