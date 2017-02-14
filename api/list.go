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
	// EmptyList represents the empty list and the terminal of a List
	EmptyList = &List{nil, nil}

	// Nil is an alias for EmptyList
	Nil = EmptyList
)

// NewList creates a new List instance
func NewList(v Value) *List {
	return &List{v, EmptyList}
}

// Cons constructs a new List by prepending to the current List
func (l *List) Cons(v Value) *List {
	return &List{v, l}
}

func (l *List) duplicate() (*List, *List) {
	if l == EmptyList {
		return l, l
	}

	h := &List{l.Value, EmptyList}
	t := h
	for li := l.Rest; li != EmptyList; li = li.Rest {
		t.Rest = &List{li.Value, EmptyList}
		t = t.Rest
	}
	return h, t
}

// Conj constructs a new List by appending to the current List
// Very expensive operation because it performs duplication
func (l *List) Conj(v Value) *List {
	h, t := l.duplicate()
	t.Rest = &List{v, EmptyList}
	return h
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
	r := l.current.Value
	l.current = l.current.Rest
	return r, true
}

// Iterable returns a new Iterable from the Iterator's current state
func (l *listIterator) Iterable() Iterable {
	return l.current
}

// Count returns the length of the List
func (l *List) Count() int {
	c := 0
	for li := l; li != EmptyList; li = li.Rest {
		c++
	}
	return c
}

// Evaluate makes a List Evaluable
func (l *List) Evaluate(c *Context) Value {
	if l == EmptyList {
		return EmptyList
	}
	if f, ok := l.Value.(*Function); ok {
		return f.Exec(c, l.Rest)
	}
	if s, ok := l.Value.(*Symbol); ok {
		if v, ok := c.Get(s.Name); ok {
			if cv, ok := v.(*Function); ok {
				return cv.Exec(c, l.Rest)
			}
		}
	}
	panic(NonFunction)
}

func (l *List) String() string {
	var b bytes.Buffer

	b.WriteString("(")
	for li := l; li != EmptyList; li = li.Rest {
		if s, ok := li.Value.(fmt.Stringer); ok {
			b.WriteString(s.String())
		} else {
			b.WriteString(li.Value.(string))
		}
		if li.Rest != EmptyList {
			b.WriteString(" ")
		}
	}
	b.WriteString(")")
	return b.String()
}
