package api

// ValueMapper returns a mapped representation of the specified Value
type ValueMapper func(Value) Value

// ValueFilter returns true if the Value remains part of a filtered Sequence
type ValueFilter func(Value) bool

// lazyMapper will only resolve its first and rest when requested
type lazyMapper struct {
	sequence Sequence
	mapper   ValueMapper
	resolved bool
	isSeq    bool
	first    Value
	rest     *lazyMapper
}

// NewLazyMapper creates a new lazy mapping Sequence that wraps the
// original and only maps its values on demand
func NewLazyMapper(s Sequence, m ValueMapper) Sequence {
	return &lazyMapper{
		sequence: s,
		mapper:   m,
	}
}

func (l *lazyMapper) resolve() *lazyMapper {
	if l.resolved {
		return l
	}

	if l.sequence.IsSequence() {
		l.isSeq = true
		l.first = l.mapper(l.sequence.First())
		l.rest = &lazyMapper{
			sequence: l.sequence.Rest(),
			mapper:   l.mapper,
		}
	}
	l.resolved = true
	l.sequence = nil
	l.mapper = nil
	return l
}

// First returns the first mapped Value from the lazyMapper
func (l *lazyMapper) First() Value {
	return l.resolve().first
}

// Rest returns the rest of the lazyMapper
func (l *lazyMapper) Rest() Sequence {
	return l.resolve().rest
}

// IsSequence returns whether this instance is a consumable Sequence
func (l *lazyMapper) IsSequence() bool {
	return l.resolve().isSeq
}

// Prepend creates a new Sequence by prepending a Value (won't be mapped)
func (l *lazyMapper) Prepend(v Value) Sequence {
	return &lazyMapper{
		first:    v,
		rest:     l,
		resolved: true,
		isSeq:    true,
	}
}

// lazyFilter will only resolve its first and rest when requested
type lazyFilter struct {
	sequence Sequence
	filter   ValueFilter
	resolved bool
	isSeq    bool
	first    Value
	rest     *lazyFilter
}

// NewLazyFilter creates a new lazy filter Sequence that wraps the
// original and only filters the next Value(s) on demand
func NewLazyFilter(s Sequence, f ValueFilter) Sequence {
	return &lazyFilter{
		sequence: s,
		filter:   f,
	}
}

func (l *lazyFilter) resolve() *lazyFilter {
	if l.resolved {
		return l
	}

	for s := l.sequence; s.IsSequence(); s = s.Rest() {
		if v := s.First(); l.filter(v) {
			l.isSeq = true
			l.first = v
			l.rest = &lazyFilter{
				sequence: s.Rest(),
				filter:   l.filter,
			}
			break
		}
	}

	l.resolved = true
	l.sequence = nil
	l.filter = nil
	return l
}

// First returns the first mapped Value from the lazyFilter
func (l *lazyFilter) First() Value {
	return l.resolve().first
}

// Rest returns the rest of the lazyFilter
func (l *lazyFilter) Rest() Sequence {
	return l.resolve().rest
}

// IsSequence returns whether this instance is a consumable Sequence
func (l *lazyFilter) IsSequence() bool {
	return l.resolve().isSeq
}

// Prepend creates a new Sequence by prepending a Value (won't be mapped)
func (l *lazyFilter) Prepend(v Value) Sequence {
	return &lazyFilter{
		first:    v,
		rest:     l,
		resolved: true,
		isSeq:    true,
	}
}
