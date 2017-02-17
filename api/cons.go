package api

import (
	"bytes"
	"fmt"
)

// NonFunction is the error returned when a non-Function is invoked
const NonFunction = "cannot resolve first element as a function"

// NonList is the error returned when a non-List is invoked
const NonList = "cannot resolve second element as a list"

// IndexNotCons is the error returns when walking a non-List
const IndexNotCons = "%d is not a Cons cell"

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
	c := i.current
	if c == Nil {
		return Nil, false
	}
	r := c.Car
	if cdr, ok := c.Cdr.(*Cons); ok {
		i.current = cdr
	} else {
		r = c
		i.current = Nil
	}
	return r, true
}

// Slice returns a new Iterable from the Iterator's current state
func (i *consIterator) Slice() Sequence {
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

// Get returns the Value at the indexed position in the List
func (c *Cons) Get(index int) Value {
	var i = 0
	for e := c; e != Nil; {
		if i == index {
			return e.Car
		}
		if n, ok := e.Cdr.(*Cons); ok {
			e = n
			i++
			continue
		}
		panic(fmt.Sprintf(IndexNotCons, index))
	}
	return Nil
}

// Eval makes a List Evaluable
func (c *Cons) Eval(ctx *Context) Value {
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
	panic(NonList)
}

func (c *Cons) listString() string {
	var b bytes.Buffer
	b.WriteString("(")
	for e := c; e != Nil; {
		b.WriteString(String(e.Car))
		if n, ok := e.Cdr.(*Cons); ok {
			e = n
			if e != Nil {
				b.WriteString(" ")
			}
		} else {
			b.WriteString(" . ")
			b.WriteString(String(e.Cdr))
			break
		}
	}
	b.WriteString(")")
	return b.String()
}

func (c *Cons) consString() string {
	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(String(c.Car))
	b.WriteString(" . ")
	b.WriteString(String(c.Cdr))
	b.WriteString(")")
	return b.String()
}

func (c *Cons) String() string {
	if _, ok := c.Cdr.(*Cons); ok {
		return c.listString()
	}
	return c.consString()
}
