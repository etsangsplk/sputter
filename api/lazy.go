package api

// ValueMapper returns a mapped representation of the specified Value
type ValueMapper func(Value) Value

// ValueFilter returns true if the Value remains part of a filtered Sequence
type ValueFilter func(Value) bool

type lazyMapper struct {
	sequence Sequence
	mapper   ValueMapper
	isSeq    bool
	first    Value
	rest     Sequence
}

// NewMapper creates a new lazy mapping Sequence that wraps the
// original and only maps its Values on demand
func NewMapper(s Sequence, m ValueMapper) Sequence {
	return &lazyMapper{
		sequence: s,
		mapper:   m,
		rest:     EmptyList,
	}
}

func (l *lazyMapper) resolve() *lazyMapper {
	if l.sequence == nil {
		return l
	}

	if l.sequence.IsSequence() {
		l.isSeq = true
		l.first = l.mapper(l.sequence.First())
		l.rest = &lazyMapper{
			sequence: l.sequence.Rest(),
			mapper:   l.mapper,
			rest:     EmptyList,
		}
	}
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
		first: v,
		rest:  l,
		isSeq: true,
	}
}

func (l *lazyMapper) Type() Name {
	return "map"
}

type lazyFilter struct {
	sequence Sequence
	filter   ValueFilter
	isSeq    bool
	first    Value
	rest     Sequence
}

// NewFilter creates a new lazy filter Sequence that wraps the
// original and only filters the next Value(s) on demand
func NewFilter(s Sequence, f ValueFilter) Sequence {
	return &lazyFilter{
		sequence: s,
		filter:   f,
		rest:     EmptyList,
	}
}

func (l *lazyFilter) resolve() *lazyFilter {
	if l.sequence == nil {
		return l
	}

	for i := l.sequence; i.IsSequence(); i = i.Rest() {
		if v := i.First(); l.filter(v) {
			l.isSeq = true
			l.first = v
			l.rest = &lazyFilter{
				sequence: i.Rest(),
				filter:   l.filter,
				rest:     EmptyList,
			}
			break
		}
	}
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
		first: v,
		rest:  l,
		isSeq: true,
	}
}

func (l *lazyFilter) Type() Name {
	return "filter"
}

type lazyConcat struct {
	sequence Sequence
	first    Value
	rest     Sequence
	isSeq    bool
}

// NewConcat creates a new sequence based on the content of
// several provided Sequences. The results are computed lazily
func NewConcat(s Sequence) Sequence {
	return &lazyConcat{
		sequence: s,
		rest:     EmptyList,
	}
}

func (l *lazyConcat) resolve() *lazyConcat {
	if l.sequence == nil {
		return l
	}

	for i := l.sequence; i.IsSequence(); i = i.Rest() {
		if f := AssertSequence(i.First()); f.IsSequence() {
			l.first = f.First()
			l.rest = &lazyConcat{
				sequence: i.Rest().Prepend(f.Rest()),
				rest:     EmptyList,
			}
			l.sequence = nil
			l.isSeq = true
			return l
		}
	}

	l.sequence = nil
	return l
}

func (l *lazyConcat) First() Value {
	return l.resolve().first
}

func (l *lazyConcat) Rest() Sequence {
	return l.resolve().rest
}

func (l *lazyConcat) IsSequence() bool {
	return l.resolve().isSeq
}

func (l *lazyConcat) Prepend(v Value) Sequence {
	return &lazyConcat{
		first: v,
		rest:  l,
		isSeq: true,
	}
}

func (l *lazyConcat) Type() Name {
	return "concat"
}
