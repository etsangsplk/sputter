package api

import (
	"sync"
	"sync/atomic"
)

// Map creates a new mapped Sequence
func Map(c Context, s Sequence, mapper Applicable) Sequence {
	var res LazyResolver
	next := s

	res = func() (Value, Sequence, bool) {
		if f, r, ok := next.Split(); ok {
			m := mapper.Apply(c, Values{f})
			next = r
			return m, NewLazySequence(res), true
		}
		return Nil, EmptyList, false
	}
	return NewLazySequence(res)
}

// MapParallel creates a new mapped Sequence from a Sequence of Sequences
// that are used to provide multiple arguments to the mapper function
func MapParallel(c Context, s Sequence, mapper Applicable) Sequence {
	var res LazyResolver
	next := make([]Sequence, 0)
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		next = append(next, f.(Sequence))
	}
	nextLen := len(next)

	res = func() (Value, Sequence, bool) {
		var exhausted int32
		args := make(Values, nextLen)

		var wg sync.WaitGroup
		wg.Add(nextLen)

		for i, s := range next {
			go func(i int, s Sequence) {
				if f, r, ok := s.Split(); ok {
					args[i] = f
					next[i] = r
				} else {
					atomic.StoreInt32(&exhausted, 1)
				}
				wg.Done()
			}(i, s)
		}
		wg.Wait()

		if exhausted > 0 {
			return Nil, EmptyList, false
		}
		m := mapper.Apply(c, args)
		return m, NewLazySequence(res), true
	}
	return NewLazySequence(res)
}

// Filter creates a new filtered Sequence
func Filter(c Context, s Sequence, filter Applicable) Sequence {
	var res LazyResolver
	next := s

	res = func() (Value, Sequence, bool) {
		for f, r, ok := next.Split(); ok; f, r, ok = r.Split() {
			next = r
			if Truthy(filter.Apply(c, Values{f})) {
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
			v := f.(Sequence)
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
func Reduce(c Context, s Sequence, reduce Applicable) Value {
	AssertMinimumArity(s, 2)
	f, r, ok := s.Split()
	res := f
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		res = reduce.Apply(c, Values{res, f})
	}
	return res
}
