package api

// ValueMapper returns a mapped representation of the specified Value
type ValueMapper func(Value) Value

// ValueFilter returns true if the Value remains part of a filtered Sequence
type ValueFilter func(Value) bool

type compElement struct {
	first Value
	rest  Sequence
}

func (c *compElement) First() Value {
	return c.first
}

func (c *compElement) Rest() Sequence {
	return c.rest
}

func (c *compElement) IsSequence() bool {
	return true
}

func (c *compElement) Prepend(v Value) Sequence {
	return &compElement{
		first: v,
		rest:  c,
	}
}

func (c *compElement) Str() Str {
	return MakeSequenceStr(c)
}

// Map creates a new mapped Sequence
func Map(s Sequence, mapper ValueMapper) Sequence {
	if !s.IsSequence() {
		return EmptyList
	}
	f := &compElement{}
	var l = f
	for i := s; i.IsSequence(); i = i.Rest() {
		n := &compElement{first: mapper(i.First())}
		l.rest = n
		l = n
	}
	l.rest = EmptyList
	return f.rest
}

// Filter creates a new filtered Sequence
func Filter(s Sequence, filter ValueFilter) Sequence {
	if !s.IsSequence() {
		return EmptyList
	}
	f := &compElement{}
	var l = f
	for i := s; i.IsSequence(); i = i.Rest() {
		v := i.First()
		if filter(v) {
			n := &compElement{first: v}
			l.rest = n
			l = n
		}
	}
	l.rest = EmptyList
	return f.rest
}

// Concat creates a new sequence based on the content of several Sequences
func Concat(s Sequence) Sequence {
	if !s.IsSequence() {
		return EmptyList
	}

	f := &compElement{}
	l := f
	for i := s; i.IsSequence(); i = i.Rest() {
		for j := AssertSequence(i.First()); j.IsSequence(); j = j.Rest() {
			n := &compElement{first: j.First()}
			l.rest = n
			l = n
		}
	}
	l.rest = EmptyList
	return f.rest
}

// Take creates a Sequence based on the first elements of the source
func Take(s Sequence, count int) Sequence {
	if !s.IsSequence() {
		return EmptyList
	}
	f := &compElement{}
	l := f
	for i, e := 0, s; i < count; i++ {
		n := &compElement{first: e.First()}
		l.rest = n
		l = n
		e = e.Rest()
	}
	l.rest = EmptyList
	return f.rest
}

// Drop creates a Sequence based on dropping the first elements of the source
func Drop(s Sequence, count int) Sequence {
	e := s
	for i := 0; i < count && e.IsSequence(); i++ {
		e = e.Rest()
	}
	return e
}

// Reduce performs a reduce operation over a Sequence, starting with the
// first two elements of that sequence.
func Reduce(c Context, s Sequence, a Applicable) Value {
	AssertMinimumArity(s, 2)
	r := s.First()
	for i := s.Rest(); i.IsSequence(); i = i.Rest() {
		r = a.Apply(c, NewVector(r, i.First()))
	}
	return r
}
