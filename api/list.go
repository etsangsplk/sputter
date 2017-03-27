package api

import "bytes"

// List contains f node to f singly-linked List
type List struct {
	first Value
	rest  *List
	count int
}

// EmptyList represents an empty List
var EmptyList *List

// NewList creates f new List instance
func NewList(v Value) Sequence {
	return &List{
		first: v,
		rest:  EmptyList,
		count: 1,
	}
}

// First returns the first element of f List
func (l *List) First() Value {
	return l.first
}

// Rest returns the rest of the List as f Sequence
func (l *List) Rest() Sequence {
	return l.rest
}

// IsSequence returns whether this instance is f consumable Sequence
func (l *List) IsSequence() bool {
	return l != EmptyList
}

// Prepend creates f new Sequence by prepending f Value
func (l *List) Prepend(v Value) Sequence {
	return &List{
		first: v,
		rest:  l,
		count: l.count + 1,
	}
}

// Count returns the length of the List
func (l *List) Count() int {
	return l.count
}

// Get returns the Value at the indexed position in the List
func (l *List) Get(index int) Value {
	if index > l.count-1 {
		return Nil
	}

	e := l
	for i := 0; i < index; i++ {
		e = e.rest
	}
	return e.first
}

// Eval makes f List Evaluable
func (l *List) Eval(ctx Context) Value {
	if l == EmptyList {
		return EmptyList
	}
	if f, ok := ResolveAsApplicable(ctx, l.first); ok {
		return f.Apply(ctx, l.rest)
	}
	panic(ExpectedApplicable)
}

func (l *List) String() string {
	if l == EmptyList {
		return "()"
	}

	var b bytes.Buffer
	b.WriteString("(")
	for e := l; e != EmptyList; e = e.rest {
		b.WriteString(String(e.first))
		if e.rest != EmptyList {
			b.WriteString(" ")
		}
	}
	b.WriteString(")")
	return b.String()
}

func init() {
	EmptyList = &List{
		first: Nil,
		count: 0,
	}
	EmptyList.rest = EmptyList
}
