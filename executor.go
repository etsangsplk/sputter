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
	var last Value
	for v := e.coder.Next(); v != EndOfCoder; v = e.coder.Next() {
		last = Evaluate(c, v)
	}
	return last
}

// Evaluate a Value against a Context
func Evaluate(c *Context, v Value) Value {
	if eval, ok := v.(Evaluable); ok {
		return eval.Evaluate(c)
	}
	return v
}
