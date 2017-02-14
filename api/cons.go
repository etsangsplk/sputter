package api

import (
	"bytes"
	"fmt"
)

// NonFunction is the error returned when a non-Function is invoked
const NonFunction = "first element of list is not a function"

// Cons contains a bound pair that can be used for constructing Lists
type Cons struct {
	Car Value
	Cdr Value
}

var (
	// Nil represents an empty Cons and the terminator of a List
	Nil = &Cons{nil, nil}
)

// NewList creates a new List instance of the form &Cons{Value, Nil}
func NewList(v Value) *Cons {
	return &Cons{v, Nil}
}

type consIterator struct {
	current *Cons
}

// Iterate creates a new Iterator instance for the Cons
func (c *Cons) Iterate() Iterator {
	return &consIterator{c}
}

// Next returns the next Value from the Iterator
func (i *consIterator) Next() (v Value, ok bool) {
	if i.current == Nil {
		return Nil, false
	}
	r := i.current.Car
	if cdr, ok := i.current.Cdr.(*Cons); ok {
		i.current = cdr
	} else {
		r = i.current
		i.current = Nil
	}
	return r, true
}

// Iterable returns a new Iterable from the Iterator's current state
func (i *consIterator) Iterable() Iterable {
	return i.current
}

// Count returns the length of the Cons
func (c *Cons) Count() int {
	i := c.Iterate()
	r := 0
	for _, ok := i.Next(); ok; _, ok = i.Next() {
		r++
	}
	return r
}

// Evaluate makes a Cons Evaluable
func (c *Cons) Evaluate(ctx *Context) Value {
	if c == Nil {
		return Nil
	}
	if a, ok := c.Cdr.(*Cons); ok {
		if f, ok := c.Car.(*Function); ok {
			return f.Exec(ctx, a)
		}

		if s, ok := c.Car.(*Symbol); ok {
			if v, ok := ctx.Get(s.Name); ok {
				if cv, ok := v.(*Function); ok {
					return cv.Exec(ctx, a)
				}
			}
		}
		panic(NonFunction)
	}
	panic("Cdr is not a Cons!")
}

func (c *Cons) listString() string {
	var b bytes.Buffer

	b.WriteString("(")
	i := c.Iterate()
	var n = false
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		if n {
			b.WriteString(" ")
		}
		if s, ok := v.(fmt.Stringer); ok {
			b.WriteString(s.String())
		} else {
			b.WriteString(v.(string))
		}
		n = true
	}
	b.WriteString(")")
	return b.String()
}

func str(v Value) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return v.(string)
}

func (c *Cons) consString() string {
	var b bytes.Buffer
	b.WriteString(str(c.Car))
	b.WriteString(" . ")
	b.WriteString(str(c.Cdr))
	return b.String()
}

func (c *Cons) String() string {
	if _, ok := c.Cdr.(*Cons); ok {
		return c.listString()
	}
	return c.consString()
}
