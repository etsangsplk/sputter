package api

import "bytes"

const (
	// ExpectedCons is thrown when a Value is not a Cons cell
	ExpectedCons = "value is not a cons cell"

	// ExpectedList is thrown when a Value is not a Cons cell
	ExpectedList = "value is not a list"

	// ExpectedSequence is thrown when a Value is not a Sequence
	ExpectedSequence = "value is not a list or vector"
)

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
			return f.Apply(ctx, a)
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

// AssertCons will cast a Value into a Cons or explode violently
func AssertCons(v Value) *Cons {
	if r, ok := v.(*Cons); ok {
		return r
	}
	panic(ExpectedCons)
}

// AssertSequence will cast a Value into a Sequence or explode violently
func AssertSequence(v Value) Sequence {
	if r, ok := v.(Sequence); ok {
		return r
	}
	panic(ExpectedSequence)
}
