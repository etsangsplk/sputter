package api

type (
	// ValueMapper returns a mapped representation of the specified Value
	ValueMapper func(Value) Value

	// ValueFilter returns true if the Value remains part of a Sequence
	ValueFilter func(Value) bool

	// ValueReducer combines two Values in some way in reducing a Sequence
	ValueReducer func(Value, Value) Value
)

// Map creates a new mapped Sequence
func Map(s Sequence, mapper ValueMapper) Sequence {
	var res LazyResolver
	e := s

	res = func() (bool, Value, Sequence) {
		if e.IsSequence() {
			r := mapper(e.First())
			e = e.Rest()
			return true, r, NewLazySequence(res)
		}
		return false, Nil, EmptyList
	}
	return NewLazySequence(res)
}

// Filter creates a new filtered Sequence
func Filter(s Sequence, filter ValueFilter) Sequence {
	var res LazyResolver
	e := s

	res = func() (bool, Value, Sequence) {
		for e.IsSequence() {
			v := e.First()
			e = e.Rest()
			if filter(v) {
				return true, v, NewLazySequence(res)
			}
		}
		return false, Nil, EmptyList
	}
	return NewLazySequence(res)
}

// Concat creates a new sequence based on the content of several Sequences
func Concat(s Sequence) Sequence {
	var res LazyResolver
	e := s

	res = func() (bool, Value, Sequence) {
		for e.IsSequence() {
			v := AssertSequence(e.First())
			e = e.Rest()
			if v.IsSequence() {
				r := v.First()
				e = e.Prepend(v.Rest())
				return true, r, NewLazySequence(res)
			}
		}
		return false, Nil, EmptyList
	}
	return NewLazySequence(res)
}

// Take creates a Sequence based on the first elements of the source
func Take(s Sequence, count int) Sequence {
	var res LazyResolver
	i, e := 0, s

	res = func() (bool, Value, Sequence) {
		if e.IsSequence() && i < count {
			r := e.First()
			e = e.Rest()
			i++
			return true, r, NewLazySequence(res)
		}
		return false, Nil, EmptyList
	}
	return NewLazySequence(res)
}

// Drop creates a Sequence based on dropping the first elements of the source
func Drop(s Sequence, count int) Sequence {
	var first, rest LazyResolver
	e := s

	first = func() (bool, Value, Sequence) {
		for i := 0; i < count && e.IsSequence(); i++ {
			e = e.Rest()
		}
		return rest()
	}

	rest = func() (bool, Value, Sequence) {
		if e.IsSequence() {
			r := e.First()
			e = e.Rest()
			return true, r, NewLazySequence(rest)
		}
		return false, Nil, EmptyList
	}

	return NewLazySequence(first)
}

// Reduce performs a reduce operation over a Sequence, starting with the
// first two elements of that sequence.
func Reduce(s Sequence, reduce ValueReducer) Value {
	AssertMinimumArity(s, 2)
	r := s.First()
	var t Value
	for i := s.Rest(); i.IsSequence(); {
		t, i = i.Split()
		r = reduce(r, t)
	}
	return r
}
