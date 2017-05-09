package parser

import (
	"regexp"

	a "github.com/kode4food/sputter/api"
)

const (
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

type mode bool

const (
	readCode mode = false
	readData mode = true
)

var keywordIdentifier = regexp.MustCompile(`^:[^(){}\[\]\s,]+`)

var specialNames = a.Variables{
	"true":  a.True,
	"false": a.False,
	"nil":   a.Nil,
}

var eofToken = &Token{
	Type:  EndOfFile,
	Value: a.Nil,
}

type endOfReader struct {
}

func (e endOfReader) Str() a.Str {
	return a.Str("")
}

// EndOfReader represents the end of a Reader stream
var EndOfReader = endOfReader{}

// Reader is responsible for returning the next Value from a Reader source
type Reader interface {
	Next() a.Value
}

// tokenReader is responsible for taking a stream of lexer Tokens and
// converting them into Lists for evaluation
type tokenReader struct {
	context a.Context
	iter    *a.Iterator
}

// NewReader instantiates a new TokenReader using the provided lexer
func NewReader(context a.Context, lexer a.Sequence) Reader {
	return &tokenReader{
		context: context,
		iter:    a.Iterate(lexer),
	}
}

func (r *tokenReader) nextToken() *Token {
	if t, ok := r.iter.Next(); ok {
		return t.(*Token)
	}
	return eofToken
}

// Next returns the next value from the Reader
func (r *tokenReader) Next() a.Value {
	return r.token(r.nextToken(), readCode)
}

func (r *tokenReader) nextData() a.Value {
	return r.token(r.nextToken(), readData)
}

func (r *tokenReader) token(t *Token, m mode) a.Value {
	switch t.Type {
	case QuoteMarker:
		return r.readQuoted()
	case ListStart:
		return r.readList(m)
	case VectorStart:
		return r.readVector(m)
	case MapStart:
		return r.readMap(m)
	case Identifier:
		return readIdentifier(t)
	case ListEnd:
		panic(UnmatchedListEnd)
	case VectorEnd:
		panic(UnmatchedVectorEnd)
	case MapEnd:
		panic(UnmatchedMapEnd)
	case EndOfFile:
		return EndOfReader
	default:
		return t.Value
	}
}

func (r *tokenReader) readQuoted() a.Quoted {
	return a.Quote(r.nextData())
}

func (r *tokenReader) readList(m mode) a.Value {
	var handle func(t *Token, m mode) a.Sequence
	var rest func(m mode) a.Sequence
	var first func() a.Value

	handle = func(t *Token, m mode) a.Sequence {
		switch t.Type {
		case ListEnd:
			return a.EmptyList
		case EndOfFile:
			panic(ListNotClosed)
		default:
			v := r.token(t, m)
			l := rest(m)
			return l.Prepend(v)
		}
	}

	rest = func(m mode) a.Sequence {
		return handle(r.nextToken(), m)
	}

	first = func() a.Value {
		t := r.nextToken()
		if f, ok := r.function(t); ok {
			if mac, ok := f.(a.Macro); ok {
				dm := mode(mac.DataMode())
				return mac.Apply(r.context, rest(dm))
			}
			return rest(m).Prepend(f)
		}
		return handle(t, m)
	}

	if m == readData {
		return rest(m)
	}
	return first()
}

func (r *tokenReader) function(t *Token) (a.Value, bool) {
	if t.Type != Identifier {
		return nil, false
	}

	if s, ok := readIdentifier(t).(a.Symbol); ok {
		if v, ok := s.Resolve(r.context); ok {
			if f, ok := v.(a.Applicable); ok {
				return f.(a.Value), true
			}
		}
	}
	return nil, false
}

func (r *tokenReader) readVector(m mode) a.Vector {
	res := make(a.Vector, 0)

	for {
		t := r.nextToken()
		switch t.Type {
		case VectorEnd:
			return res
		case EndOfFile:
			panic(VectorNotClosed)
		default:
			e := r.token(t, m)
			res = append(res, e)
		}
	}
}

func (r *tokenReader) readMap(m mode) a.Associative {
	res := make(a.Associative, 0)
	mp := make(a.Vector, 2)

	for idx := 0; ; idx++ {
		t := r.nextToken()
		switch t.Type {
		case MapEnd:
			if idx%2 == 0 {
				return res
			}
			panic(MapNotPaired)
		case EndOfFile:
			panic(MapNotClosed)
		default:
			e := r.token(t, m)
			if idx%2 == 0 {
				mp[0] = e
			} else {
				mp[1] = e
				res = append(res, mp)
				mp = make(a.Vector, 2)
			}
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

// EvalReader consumes a Reader and evaluates the resulting Sequence
func EvalReader(c a.Context, reader Reader) a.Value {
	s := a.Vector{}
	for v := reader.Next(); v != EndOfReader; v = reader.Next() {
		s = s.Conjoin(v).(a.Vector)
	}
	return a.EvalSequence(c, s)
}
