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
	next := s

	res = func() (Value, Sequence, bool) {
		if f, r, ok := next.Split(); ok {
			m := mapper(f)
			next = r
			return m, NewLazySequence(res), true
		}
		return Nil, EmptyList, false
	}
	return NewLazySequence(res)
}

// Filter creates a new filtered Sequence
func Filter(s Sequence, filter ValueFilter) Sequence {
	var res LazyResolver
	next := s

	res = func() (Value, Sequence, bool) {
		for f, r, ok := next.Split(); ok; f, r, ok = r.Split() {
			next = r
			if filter(f) {
				return f, NewLazySequence(res), true
			}
		}
		return Nil, EmptyList, false
	}
	return NewLazySequence(res)
}

// Concat creates a new sequence based on the content of several Sequences
func Concat(s Sequence) Sequence {
	var res LazyResolver
	next := s

	res = func() (Value, Sequence, bool) {
		for f, r, ok := next.Split(); ok; f, r, ok = r.Split() {
			v := AssertSequence(f)
			next = r
			if vf, vr, ok := v.Split(); ok {
				next = next.Prepend(vr)
				return vf, NewLazySequence(res), true
			}
		}
		return Nil, EmptyList, false
	}
	return NewLazySequence(res)
}

// Take creates a Sequence based on the first elements of the source
func Take(s Sequence, count int) Sequence {
	var res LazyResolver
	next := s
	idx := 0

	res = func() (Value, Sequence, bool) {
		if f, r, ok := next.Split(); ok && idx < count {
			next = r
			idx++
			return f, NewLazySequence(res), true
		}
		return Nil, EmptyList, false
	}
	return NewLazySequence(res)
}

// Drop creates a Sequence based on dropping the first elements of the source
func Drop(s Sequence, count int) Sequence {
	var first, rest LazyResolver
	next := s

	first = func() (Value, Sequence, bool) {
		for i := 0; i < count && next.IsSequence(); i++ {
			next = next.Rest()
		}
		return rest()
	}

	rest = func() (Value, Sequence, bool) {
		if f, r, ok := next.Split(); ok {
			next = r
			return f, NewLazySequence(rest), true
		}
		return Nil, EmptyList, false
	}

	return NewLazySequence(first)
}

// Reduce performs a reduce operation over a Sequence, starting with the
// first two elements of that sequence.
func Reduce(s Sequence, reduce ValueReducer) Value {
	AssertMinimumArity(s, 2)
	f, r, ok := s.Split()
	res := f
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		res = reduce(res, f)
	}
	return res
}
