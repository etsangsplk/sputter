package api

// ValueMapper returns a mapped representation of the specified Value
type ValueMapper func(Value) Value

// lazySequence will only resolve its first and rest when requested
type lazySequence struct {
	sequence Sequence
	mapper   ValueMapper
	resolved bool
	isSeq    bool
	first    Value
	rest     *lazySequence
}

// NewLazySequence creates a new LazySequence that wraps the original
// and only maps its values on demand
func NewLazySequence(s Sequence, m ValueMapper) Sequence {
	return &lazySequence{
		sequence: s,
		mapper:   m,
	}
}

func (l *lazySequence) resolve() *lazySequence {
	if l.resolved {
		return l
	}

	l.isSeq = l.sequence.IsSequence()
	if l.isSeq {
		l.first = l.mapper(l.sequence.First())
		l.rest = &lazySequence{
			sequence: l.sequence.Rest(),
			mapper:   l.mapper,
		}
	}
	l.resolved = true
	l.sequence = nil
	l.mapper = nil
	return l
}

// First returns the first mapped Value from the LazySequence
func (l *lazySequence) First() Value {
	return l.resolve().first
}

// Rest returns the rest of the LazySequence
func (l *lazySequence) Rest() Sequence {
	return l.resolve().rest
}

// IsSequence returns whether this instance is a consumable Sequence
func (l *lazySequence) IsSequence() bool {
	return l.resolve().isSeq
}

// Prepend creates a new Sequence by prepending a Value (won't be mapped)
func (l *lazySequence) Prepend(v Value) Sequence {
	return &lazySequence{
		first:    v,
		rest:     l,
		resolved: true,
		isSeq:    true,
	}
}
