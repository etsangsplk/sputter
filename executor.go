package sputter

import "fmt"

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
	if eval, ok := v.(Evaluable); ok {
		return eval.Evaluate(c)
	}
	return v
}

func EvaluateToString(c *Context, v Value) string {
	result := Evaluate(c, v)
	if str, ok := result.(fmt.Stringer); ok {
		return str.String()
	}
	return result.(string)
}
