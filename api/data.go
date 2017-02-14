package api

import "fmt"

// Data identifies a Value as being in data mode (literal)
type Data struct {
	Value Value
}

// Evaluate makes Data Evaluable
func (l *Data) Evaluate(c *Context) Value {
	return l.Value
}

func (l *Data) String() string {
	if s, ok := l.Value.(fmt.Stringer); ok {
		return s.String()
	}
	return l.Value.(string)
}
