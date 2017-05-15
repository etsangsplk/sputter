package parser

import (
	"regexp"

	a "github.com/kode4food/sputter/api"
)

const (
	// QuoteNotPaired is thrown when a Quote is not completed
	QuoteNotPaired = "quote without data to be quoted"

	// ListNotClosed is thrown when EOF is reached inside a List
	ListNotClosed = "end of file reached with open list"

	// UnmatchedListEnd is thrown if a list is ended without being started
	UnmatchedListEnd = "encountered ')' with no open list"

	// VectorNotClosed is thrown when EOF is reached inside a Vector
	VectorNotClosed = "end of file reached with open vector"

	// UnmatchedVectorEnd is thrown if a vector is ended without being started
	UnmatchedVectorEnd = "encountered ']' with no open vector"

	// MapNotClosed is thrown when EOF is reached inside a Map
	MapNotClosed = "end of file reached with open map"

	// UnmatchedMapEnd is thrown if a list is ended without being started
	UnmatchedMapEnd = "encountered '}' with no open map"

	// MapNotPaired is thrown if a map doesn't have an even number of elements
	MapNotPaired = "map does not contain an even number of elements"
)

var keywordIdentifier = regexp.MustCompile(`^:[^(){}\[\]\s,]+`)

var specialNames = a.Variables{
	"true":  a.True,
	"false": a.False,
	"nil":   a.Nil,
}

type reader struct {
	context a.Context
	iter    *a.Iterator
}

// NewReader completely consumes a Lexer before returning a Value Sequence
func NewReader(context a.Context, lexer a.Sequence) a.Sequence {
	r := a.NewVector()
	iter := a.Iterate(lexer)
	ri := &reader{
		context: context,
		iter:    iter,
	}
	for f, ok := ri.nextCode(); ok; f, ok = ri.nextCode() {
		r = r.Conjoin(f).(a.Vector)
	}
	return r
}

func (r *reader) nextToken() (*Token, bool) {
	if t, ok := r.iter.Next(); ok {
		return t.(*Token), true
	}
	return nil, false
}

func (r *reader) nextCode() (a.Value, bool) {
	return r.nextValue(false)
}

func (r *reader) nextData() (a.Value, bool) {
	return r.nextValue(true)
}

func (r *reader) nextValue(data bool) (a.Value, bool) {
	if t, ok := r.nextToken(); ok {
		return r.value(t, data), true
	}
	return a.Nil, false
}

func (r *reader) value(t *Token, data bool) a.Value {
	var v a.Value
	switch t.Type {
	case QuoteMarker:
		v = r.quoted(data)
	case ListStart:
		v = r.list(data)
	case VectorStart:
		v = r.vector(data)
	case MapStart:
		v = r.associative(data)
	case Identifier:
		v = readIdentifier(t, data)
	case ListEnd:
		panic(UnmatchedListEnd)
	case VectorEnd:
		panic(UnmatchedVectorEnd)
	case MapEnd:
		panic(UnmatchedMapEnd)
	default:
		v = t.Value
	}
	return r.expand(v)
}

func (r *reader) quoted(data bool) a.Value {
	if v, ok := r.nextData(); ok {
		q := a.NewQualifiedSymbol("quote", a.BuiltInDomain)
		l := a.NewList(v).Prepend(q)
		if data {
			return l
		}
		return l.(a.List).Evaluable()
	}
	panic(QuoteNotPaired)
}

func (r *reader) list(data bool) a.Value {
	var handle func(t *Token) a.List
	var rest func() a.List
	var first func() a.List

	handle = func(t *Token) a.List {
		switch t.Type {
		case ListEnd:
			return a.EmptyList
		default:
			v := r.value(t, data)
			l := rest()
			return l.Prepend(v).(a.List)
		}
	}

	rest = func() a.List {
		if t, ok := r.nextToken(); ok {
			return handle(t)
		}
		panic(ListNotClosed)
	}

	first = func() a.List {
		if t, ok := r.nextToken(); ok {
			if t.Type != Identifier {
				return handle(t)
			}
			v := r.value(t, false)
			_, data = r.macro(v)
			return rest().Prepend(v).(a.List)
		}
		panic(ListNotClosed)
	}

	if data {
		return rest()
	}
	return first().Evaluable()
}

func (r *reader) expand(v a.Value) a.Value {
	if l, ok := v.(a.List); ok {
		if _, ok := l.(a.Evaluable); ok {
			if m, ok := r.macro(l.First()); ok {
				return r.expand(m.Apply(r.context, l.Rest()))
			}
		}
	}
	return v
}

func (r *reader) macro(v a.Value) (a.Function, bool) {
	if s, ok := v.(a.Symbol); ok {
		if r, ok := s.Resolve(r.context); ok {
			if f, ok := r.(a.Function); ok {
				if a.IsMacro(f) {
					return f, true
				}
			}
		}
	}
	return nil, false
}

func (r *reader) vector(data bool) a.Value {
	res := make([]a.Value, 0)

	for {
		if t, ok := r.nextToken(); ok {
			switch t.Type {
			case VectorEnd:
				v := a.NewVector(res...)
				if data {
					return v
				}
				return v.Evaluable()
			default:
				e := r.value(t, data)
				res = append(res, e)
			}
		} else {
			panic(VectorNotClosed)
		}
	}
}

func (r *reader) associative(data bool) a.Value {
	res := make([]a.Vector, 0)
	mp := make([]a.Value, 2)

	for idx := 0; ; idx++ {
		if t, ok := r.nextToken(); ok {
			switch t.Type {
			case MapEnd:
				if idx%2 == 0 {
					as := a.NewAssociative(res...)
					if data {
						return as
					}
					return as.Evaluable()
				}
				panic(MapNotPaired)
			default:
				e := r.value(t, data)
				if idx%2 == 0 {
					mp[0] = e
				} else {
					mp[1] = e
					res = append(res, a.NewVector(mp...))
					mp = make([]a.Value, 2)
				}
			}
		} else {
			panic(MapNotClosed)
		}
	}
}

func readIdentifier(t *Token, data bool) a.Value {
	n := a.Name(t.Value.(a.Str))
	if v, ok := specialNames[n]; ok {
		return v
	}

	s := string(n)
	if keywordIdentifier.MatchString(s) {
		return a.NewKeyword(n[1:])
	}
	sym := a.ParseSymbol(n)
	if data {
		return sym
	}
	return sym.Evaluable()
}
