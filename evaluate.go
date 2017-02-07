package main

// Evaluate a Value against a Context
func Evaluate(c *Context, v Value) Value {
	if eval, ok := v.(Evaluable); ok {
		return eval.Evaluate(c)
	}
	return v
}

// EvaluateCoder evaluates each element of the provided Coder
func EvaluateCoder(c *Context, coder *Coder) Value {
	var lastEval Value
	for v := coder.Next(); v != EndOfCoder; v = coder.Next() {
		lastEval = Evaluate(c, v)
	}
	return lastEval
}

// EvaluateIterator evaluates each element of the provided Iterator
func EvaluateIterator(c *Context, iter Iterator) Value {
	var lastEval Value = EmptyList
	for val, ok := iter.Next(); ok; val, ok = iter.Next() {
		lastEval = Evaluate(c, val)
	}
	return lastEval
}
