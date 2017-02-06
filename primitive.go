package sputter

import "fmt"

// Value is the generic interface for all 'Values' in the VM
type Value interface {
}

// Iterable values can be used in loops and comprehensions
type Iterable interface {
	Iterate() Iterator
}

// Iterator interfaces are stateful iteration interfaces
type Iterator interface {
	Next() (Value, bool)
	Iterable() Iterable
}

// Literal identifies a Value as being a literal reference
type Literal struct {
	value Value
}

// Evaluate makes a Literal Evaluable
func (l *Literal) Evaluate(c *Context) Value {
	return l.value
}

func (l *Literal) String() string {
	if str, ok := l.value.(fmt.Stringer); ok {
		return str.String()
	}
	return l.value.(string)
}
