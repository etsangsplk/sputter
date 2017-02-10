package api

import (
	"bytes"
	"fmt"
)

// NonFunction is the error returned when a non-Function is invoked
const NonFunction = "first element of list is not a function"

// List is a Value that maintains a singly-linked list of Values
type List struct {
	Value Value
	Rest  *List
}

var (
	// EmptyList represents the empty list and the terminal 'rest' of a List
	EmptyList = &List{nil, nil}

	// Nil is an alias for EmptyList
	Nil = EmptyList
)

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

	first := &List{l.Value, EmptyList}
	last := first
	for current := l.Rest; current != EmptyList; current = current.Rest {
		last.Rest = &List{current.Value, EmptyList}
		last = last.Rest
	}
	return first, last
}

// Conj constructs a new List by appending to the current List
// Very expensive operation because it performs duplication
func (l *List) Conj(tail Value) *List {
	first, last := l.duplicate()
	last.Rest = &List{tail, EmptyList}
	return first
}

// ListIterator is an Iterator implementation for the List type
type listIterator struct {
	current *List
}

// Iterate creates a new Iterator instance for the List
func (l *List) Iterate() Iterator {
	return &listIterator{l}
}

// Next returns the next Value from the Iterator
func (l *listIterator) Next() (Value, bool) {
	if l.current == EmptyList {
		return EmptyList, false
	}
	result := l.current.Value
	l.current = l.current.Rest
	return result, true
}

// Iterable returns a new Iterable from the Iterator's current state
func (l *listIterator) Iterable() Iterable {
	return l.current
}

// Count returns the length of the List
func (l *List) Count() int {
	count := 0
	for current := l; current != EmptyList; current = current.Rest {
		count++
	}
	return count
}

// Evaluate makes a List Evaluable
func (l *List) Evaluate(c *Context) Value {
	if l == EmptyList {
		return EmptyList
	}
	if function, ok := l.Value.(*Function); ok {
		return function.Exec(c, l.Rest)
	}
	if sym, ok := l.Value.(*Symbol); ok {
		if v, ok := c.Get(sym.Name); ok {
			if entry, ok := v.(*Function); ok {
				return entry.Exec(c, l.Rest)
			}
		}
	}
	panic(NonFunction)
}

func (l *List) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(")
	for current := l; current != EmptyList; current = current.Rest {
		if str, ok := current.Value.(fmt.Stringer); ok {
			buffer.WriteString(str.String())
		} else {
			buffer.WriteString(current.Value.(string))
		}
		if current.Rest != EmptyList {
			buffer.WriteString(" ")
		}
	}
	buffer.WriteString(")")
	return buffer.String()
}
