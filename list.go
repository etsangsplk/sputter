package main

import (
	"bytes"
	"fmt"
)

// List is a Value that maintains a singly-linked list of Values
type List struct {
	Value Value
	Rest  *List
}

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

	first := &List{l.Value, EmptyList}
	last := first
	for current := l.Rest; current != EmptyList; current = current.Rest {
		last.Rest = &List{current.Value, EmptyList}
		last = last.Rest
	}
	return first, last
}

// Conj constructs a new List by appending to the current List
func (l *List) Conj(tail Value) *List {
	first, last := l.duplicate()
	last.Rest = &List{tail, EmptyList}
	return first
}

// ListIterator is an Iterator implementation for the List type
type ListIterator struct {
	current *List
}

// Iterate creates a new Iterator instance for the List
func (l *List) Iterate() Iterator {
	return &ListIterator{l}
}

// Next returns the next Value from the Iterator
func (l *ListIterator) Next() (Value, bool) {
	if l.current == EmptyList {
		return EmptyList, false
	}
	result := l.current.Value
	l.current = l.current.Rest
	return result, true
}

// Iterable returns a new Iterable from the Iterator's current state
func (l *ListIterator) Iterable() Iterable {
	return l.current
}

// Evaluate makes a List Evaluable
func (l *List) Evaluate(c *Context) Value {
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
