package sputter

import "fmt"

// Value is the generic interface for all 'Values' in the VM
type Value interface {
}

// Iterable values can be used in loops and comprehensions
type Iterable interface {
	Iterate() Iterator
}

// Iterator functions are ones that return the next Value of an Iterable
type Iterator func() (Value, bool)

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
