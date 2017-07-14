package api

// NewRange instantiates a new lazy Range Sequence
func NewRange(first, last, step Number) Sequence {
	var inc LazyResolver
	val := first
	cmp := LessThan

	if step.Cmp(Zero) == LessThan {
		cmp = GreaterThan
	}

	inc = func() (bool, Value, Sequence) {
		if val.Cmp(last) == cmp {
			r := val
			val = val.Add(step)
			return true, r, NewLazySequence(inc)
		}
		return false, Nil, EmptyList
	}

	return NewLazySequence(inc)
}
