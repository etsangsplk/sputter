package api

type (
	// LazyResolver is used to resolve the elements of a lazy Sequence
	LazyResolver func() (bool, Value, Sequence)

	lazySequence struct {
		once     Do
		resolver LazyResolver

		isSeq  bool
		result Value
		rest   Sequence
	}
)

// NewLazySequence creates a new lazy Sequence based on the provided LazyResolver
func NewLazySequence(r LazyResolver) Sequence {
	return &lazySequence{
		once:     Once(),
		resolver: r,
		result:   Nil,
		rest:     EmptyList,
	}
}

func (l *lazySequence) resolve() *lazySequence {
	l.once(func() {
		l.isSeq, l.result, l.rest = l.resolver()
		l.resolver = nil
	})
	return l
}

func (l *lazySequence) IsSequence() bool {
	return l.resolve().isSeq
}

func (l *lazySequence) First() Value {
	return l.resolve().result
}

func (l *lazySequence) Rest() Sequence {
	return l.resolve().rest
}

func (l *lazySequence) Prepend(v Value) Sequence {
	return &lazySequence{
		once:   Never(),
		isSeq:  true,
		result: v,
		rest:   l,
	}
}

func (l *lazySequence) Type() Name {
	return "lazy-sequence"
}

func (l *lazySequence) Str() Str {
	return MakeDumpStr(l)
}
