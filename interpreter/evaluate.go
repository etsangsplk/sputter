package interpreter

import (
	a "github.com/kode4food/sputter/api"
)

// Evaluate a Value against a Context
func Evaluate(c *a.Context, v a.Value) a.Value {
	if eval, ok := v.(a.Evaluable); ok {
		return eval.Evaluate(c)
	}
	return v
}

// EvaluateCoder evaluates each element of the provided Coder
func EvaluateCoder(c *a.Context, coder *Coder) a.Value {
	var lastEval a.Value
	for v := coder.Next(); v != EndOfCoder; v = coder.Next() {
		lastEval = Evaluate(c, v)
	}
	return lastEval
}

// EvaluateIterator evaluates each element of the provided Iterator
func EvaluateIterator(c *a.Context, iter a.Iterator) a.Value {
	var lastEval a.Value = a.EmptyList
	for val, ok := iter.Next(); ok; val, ok = iter.Next() {
		lastEval = Evaluate(c, val)
	}
	return lastEval
}
