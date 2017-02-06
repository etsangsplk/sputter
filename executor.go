package sputter

// Executor is responsible for taking Code Lists and executing them against
// a Context instance
type Executor struct {
	coder *Coder
}

// NewExecutor instantiates a new Executor
func NewExecutor(c *Coder) *Executor {
	return &Executor{c}
}

// Exec invokes an Executor
func (e *Executor) Exec(c *Context) Value {
	var lastEval Value
	for v := e.coder.Next(); v != EndOfCoder; v = e.coder.Next() {
		lastEval = Evaluate(c, v)
	}
	return lastEval
}

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
