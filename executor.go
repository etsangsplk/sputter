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
	v := e.coder.Next()
	return Evaluate(c, v)
}

// Evaluate a Value against a Context
func Evaluate(c *Context, v Value) Value {
	if list, ok := v.(*List); ok {
		if function, ok := list.value.(*Function); ok {
			return function.exec(c, list.rest)
		}
	}
	return v
}
