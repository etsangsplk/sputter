package api

import "bytes"

// Nil represents an empty Cons and the terminator of a List
var Nil = &Cons{nil, nil}

// Cons contains a bound pair that can be used for constructing Lists
type Cons struct {
	Car Value
	Cdr Value
}

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
func (i *consIterator) Next() (Value, bool) {
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

// Rest returns a new Iterable from the Iterator's current state
func (i *consIterator) Rest() Sequence {
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
		panic(ExpectedCons)
	}
	return Nil
}

// Eval makes a List Evaluable
func (c *Cons) Eval(ctx Context) Value {
	if c == Nil {
		return Nil
	}

	if a, ok := c.Cdr.(*Cons); ok {
		if f, ok := ResolveAsFunction(ctx, c.Car); ok {
			return f.Exec(ctx, a)
		}
		panic(ExpectedFunction)
	}
	panic(ExpectedList)
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
	if c == Nil {
		return "()"
	}
	if _, ok := c.Cdr.(*Cons); ok {
		return c.listString()
	}
	return c.consString()
}
