package api

type (
	// LazyResolver is used to resolve the elements of a lazy Sequence
	LazyResolver func() (Value, bool, LazyResolver)

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
		v, ok, r := l.resolver()
		l.isSeq = ok
		l.result = v
		l.resolver = nil
		if ok {
			l.rest = &lazySequence{
				once:     Once(),
				resolver: r,
				result:   Nil,
				rest:     EmptyList,
			}
		}
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
