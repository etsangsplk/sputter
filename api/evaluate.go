package api

// Evaluate a Value against a Context
func Evaluate(c *Context, v Value) Value {
	if eval, ok := v.(Evaluable); ok {
		return eval.Evaluate(c)
	}
	return v
}

// EvaluateIterator evaluates each element of the provided Iterator
func EvaluateIterator(c *Context, iter Iterator) Value {
	var lastEval Value = EmptyList
	for val, ok := iter.Next(); ok; val, ok = iter.Next() {
		lastEval = Evaluate(c, val)
	}
	return lastEval
}
