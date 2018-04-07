package reader

import (
	"regexp"

	a "github.com/kode4food/sputter/api"
)

const (
	// PrefixedNotPaired is thrown when a Quote is not completed
	PrefixedNotPaired = "end of file reached before completing %s"

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

type reader struct {
	iter *a.Iterator
}

var (
	keywordIdentifier = regexp.MustCompile(`^:[^(){}\[\]\s,]+`)

	quoteSym    = a.NewBuiltInSymbol("quote")
	syntaxSym   = a.NewBuiltInSymbol("syntax-quote")
	unquoteSym  = a.NewBuiltInSymbol("unquote")
	splicingSym = a.NewBuiltInSymbol("unquote-splicing")

	specialNames = a.Variables{
		"true":  a.True,
		"false": a.False,
		"nil":   a.Nil,
	}
)

// ReadStr converts the raw source into unexpanded data structures
func ReadStr(src a.Str) a.Sequence {
	l := Scan(src)
	return Read(l)
}

// Read returns a Lazy Sequence of scanned data structures
func Read(lexer a.Sequence) a.Sequence {
	var res a.LazyResolver
	r := newReader(lexer)

	res = func() (a.Value, a.Sequence, bool) {
		if f, ok := r.nextValue(); ok {
			return f, a.NewLazySequence(res), true
		}
		return a.Nil, a.EmptyList, false
	}

	return a.NewLazySequence(res)
}

func newReader(lexer a.Sequence) *reader {
	return &reader{
		iter: a.Iterate(lexer),
	}
}

func (r *reader) nextToken() (*Token, bool) {
	if t, ok := r.iter.Next(); ok {
		return t.(*Token), true
	}
	return nil, false
}

func (r *reader) nextValue() (a.Value, bool) {
	if t, ok := r.nextToken(); ok {
		return r.value(t), true
	}
	return a.Nil, false
}

func (r *reader) value(t *Token) a.Value {
	switch t.Type {
	case QuoteMarker:
		return r.prefixed(quoteSym)
	case SyntaxMarker:
		return r.prefixed(syntaxSym)
	case UnquoteMarker:
		return r.prefixed(unquoteSym)
	case SpliceMarker:
		return r.prefixed(splicingSym)
	case ListStart:
		return r.list()
	case VectorStart:
		return r.vector()
	case MapStart:
		return r.associative()
	case Identifier:
		return readIdentifier(t)
	case ListEnd:
		panic(a.ErrStr(UnmatchedListEnd))
	case VectorEnd:
		panic(a.ErrStr(UnmatchedVectorEnd))
	case MapEnd:
		panic(a.ErrStr(UnmatchedMapEnd))
	default:
		return t.Value
	}
}

func (r *reader) prefixed(s a.Symbol) a.Value {
	if v, ok := r.nextValue(); ok {
		return a.NewList(s, v)
	}
	panic(a.ErrStr(PrefixedNotPaired, s))
}

func (r *reader) list() a.Value {
	var handle func(t *Token) *a.List
	var rest func() *a.List

	handle = func(t *Token) *a.List {
		switch t.Type {
		case ListEnd:
			return a.EmptyList
		default:
			v := r.value(t)
			l := rest()
			return l.Prepend(v).(*a.List)
		}
	}

	rest = func() *a.List {
		if t, ok := r.nextToken(); ok {
			return handle(t)
		}
		panic(a.ErrStr(ListNotClosed))
	}

	return rest()
}

func (r *reader) vector() a.Value {
	res := make(a.Vector, 0)

	for {
		if t, ok := r.nextToken(); ok {
			switch t.Type {
			case VectorEnd:
				return res
			default:
				e := r.value(t)
				res = append(res, e)
			}
		} else {
			panic(a.ErrStr(VectorNotClosed))
		}
	}
}

func (r *reader) associative() a.Value {
	res := make([]a.Vector, 0)
	mp := make(a.Vector, 2)

	for idx := 0; ; idx++ {
		if t, ok := r.nextToken(); ok {
			switch t.Type {
			case MapEnd:
				if idx%2 == 0 {
					return a.NewAssociative(res...)
				}
				panic(a.ErrStr(MapNotPaired))
			default:
				e := r.value(t)
				if idx%2 == 0 {
					mp[0] = e
				} else {
					mp[1] = e
					res = append(res, mp)
					mp = make(a.Vector, 2)
				}
			}
		} else {
			panic(a.ErrStr(MapNotClosed))
		}
	}
}

func readIdentifier(t *Token) a.Value {
	n := a.Name(t.Value.(a.Str))
	if v, ok := specialNames[n]; ok {
		return v
	}

	s := string(n)
	if keywordIdentifier.MatchString(s) {
		return a.NewKeyword(n[1:])
	}
	return a.ParseSymbol(n)
}
