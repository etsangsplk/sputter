package api

import "fmt"

// Literal identifies a Value as being a literal reference
type Literal struct {
	Value Value
}

// Evaluate makes a Literal Evaluable
func (l *Literal) Evaluate(c *Context) Value {
	return l.Value
}

func (l *Literal) String() string {
	if str, ok := l.Value.(fmt.Stringer); ok {
		return str.String()
	}
	return l.Value.(string)
}
