package sputter

import (
	"bytes"
	"fmt"
)

// NonFunction is the error returned when a non-Function is invoked
const NonFunction = "first element of list is not a function"

// List is a Value that maintains a singly-linked list of Values
type List struct {
	value Value
	rest  *List
}

// ListProcessor is the standard signature for a function that processes lists
type ListProcessor func(*Context, *List) Value

// EmptyList represents the empty list and the terminal 'rest' of a List
var EmptyList = &List{nil, nil}

// NewList creates a new List instance
func NewList(head Value) *List {
	return &List{head, EmptyList}
}

// Cons constructs a new List by prepending to the current List
func (l *List) Cons(head Value) *List {
	return &List{head, l}
}

func (l *List) duplicate() (*List, *List) {
	if l == EmptyList {
		return l, l
	}

	first := &List{l.value, EmptyList}
	last := first
	for current := l.rest; current != EmptyList; current = current.rest {
		last.rest = &List{current.value, EmptyList}
		last = last.rest
	}
	return first, last
}

// Conj constructs a new List by appending to the current List
func (l *List) Conj(tail Value) *List {
	first, last := l.duplicate()
	last.rest = &List{tail, EmptyList}
	return first
}

// Iterate creates a new Iterator instance for the List
func (l *List) Iterate() Iterator {
	current := l
	return func() (Value, bool) {
		if current == EmptyList {
			return nil, false
		}
		result := current.value
		current = current.rest
		return result, true
	}
}

func (l *List) Evaluate(c *Context) Value {
	if function, ok := l.value.(*Function); ok {
		return function.exec(c, l.rest)
	}
	return l
}

func (l *List) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(")
	for current := l; current != EmptyList; current = current.rest {
		if str, ok := current.value.(fmt.Stringer); ok {
			buffer.WriteString(str.String())
		} else {
			buffer.WriteString(current.value.(string))
		}
		if current.rest != EmptyList {
			buffer.WriteString(" ")
		}		
	}
	buffer.WriteString(")")
	return buffer.String()
}


// Function is a Value that can be invoked
type Function struct {
	name string
	exec ListProcessor
}

func (f *Function) String() string {
	return f.name
}
