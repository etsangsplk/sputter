package api

// NewRange instantiates a new lazy Range Sequence
func NewRange(first, last, step Number) Sequence {
	var inc LazyResolver
	val := first
	cmp := LessThan

	if step.Cmp(Zero) == LessThan {
		cmp = GreaterThan
	}

	inc = func() (Value, Sequence, bool) {
		if val.Cmp(last) == cmp {
			r := val
			val = val.Add(step)
			return r, NewLazySequence(inc), true
		}
		return Nil, EmptyList, false
	}

	return NewLazySequence(inc)
}
