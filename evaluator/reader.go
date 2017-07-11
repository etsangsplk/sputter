package evaluator

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
	context a.Context
	iter    *a.Iterator
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

// Read completely consumes a Lexer before returning a Value Sequence
func Read(context a.Context, lexer a.Sequence) a.Sequence {
	r := a.NewVector()
	iter := a.Iterate(lexer)
	ri := &reader{
		context: context,
		iter:    iter,
	}
	for f, ok := ri.nextValue(); ok; f, ok = ri.nextValue() {
		r = r.Conjoin(f).(a.Vector)
	}
	return a.NewBlock(r)
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
		panic(UnmatchedListEnd)
	case VectorEnd:
		panic(UnmatchedVectorEnd)
	case MapEnd:
		panic(UnmatchedMapEnd)
	default:
		return t.Value
	}
}

func (r *reader) prefixed(s a.Symbol) a.Value {
	if v, ok := r.nextValue(); ok {
		return a.NewList(s, v)
	}
	panic(a.Err(PrefixedNotPaired, s))
}

func (r *reader) list() a.Value {
	var handle func(t *Token) a.List
	var rest func() a.List

	handle = func(t *Token) a.List {
		switch t.Type {
		case ListEnd:
			return a.EmptyList
		default:
			v := r.value(t)
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

	return rest()
}

func (r *reader) vector() a.Value {
	res := make([]a.Value, 0)

	for {
		if t, ok := r.nextToken(); ok {
			switch t.Type {
			case VectorEnd:
				return a.NewVector(res...)
			default:
				e := r.value(t)
				res = append(res, e)
			}
		} else {
			panic(VectorNotClosed)
		}
	}
}

func (r *reader) associative() a.Value {
	res := make([]a.Vector, 0)
	mp := make([]a.Value, 2)

	for idx := 0; ; idx++ {
		if t, ok := r.nextToken(); ok {
			switch t.Type {
			case MapEnd:
				if idx%2 == 0 {
					return a.NewAssociative(res...)
				}
				panic(MapNotPaired)
			default:
				e := r.value(t)
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
