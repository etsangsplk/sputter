package reader

import a "github.com/kode4food/sputter/api"

const (
	// ListNotClosed is thrown when EOF is reached inside a Cons
	ListNotClosed = "end of file reached with open list"

	// UnmatchedListEnd is thrown if a list is ended without being started
	UnmatchedListEnd = "encountered ')' with no open list"

	// VectorNotClosed is thrown when EOF is reached inside a Vector
	VectorNotClosed = "end of file reached with open vector"

	// UnmatchedVectorEnd is thrown if a vector is ended without being started
	UnmatchedVectorEnd = "encountered ']' with no open vector"
)

var specialNames = a.Variables{
	"true":  a.True,
	"false": a.False,
	"nil":   a.Nil,
}

// EndOfReader represents the end of a Reader stream
var EndOfReader = struct{}{}

// Reader is responsible for returning the next Value from a Reader source
type Reader interface {
	Next() a.Value
}

// tokenReader is responsible for taking a stream of Lexer Tokens and
// converting them into Lists for evaluation
type tokenReader struct {
	context a.Context
	lexer   Lexer
}

// NewReader instantiates a new TokenReader using the provided Lexer
func NewReader(context a.Context, l Lexer) Reader {
	return &tokenReader{context, l}
}

// Next returns the next value from the Reader
func (r *tokenReader) Next() a.Value {
	return r.token(r.lexer.Next(), false)
}

func (r *tokenReader) nextData() a.Value {
	return r.token(r.lexer.Next(), true)
}

func (r *tokenReader) token(t *Token, data bool) a.Value {
	switch t.Type {
	case QuoteMarker:
		return r.quote()
	case ListStart:
		return r.list(data)
	case VectorStart:
		return r.vector(data)
	case Identifier:
		n := a.Name(t.Value.(string))
		if v, ok := specialNames[n]; ok {
			return v
		}
		return a.NewSymbol(n)
	case ListEnd:
		panic(UnmatchedListEnd)
	case VectorEnd:
		panic(UnmatchedVectorEnd)
	case EndOfFile:
		return EndOfReader
	default:
		return t.Value
	}
}

func (r *tokenReader) quote() *a.Quote {
	return &a.Quote{Value: r.nextData()}
}

func (r *tokenReader) prepareList(l *a.Cons) a.Value {
	if s, ok := l.Car.(*a.Symbol); ok {
		if v, ok := r.context.Get(s.Name); ok {
			if f, ok := v.(*a.Function); ok {
				fl := &a.Cons{Car: f, Cdr: l.Cdr}
				if f.Prepare != nil {
					return f.Prepare(r.context, fl)
				}
				return fl
			}
		}
	}
	return l
}

func (r *tokenReader) list(data bool) a.Value {
	var next func() *a.Cons

	next = func() *a.Cons {
		t := r.lexer.Next()
		switch t.Type {
		case ListEnd:
			return a.Nil
		case EndOfFile:
			panic(ListNotClosed)
		default:
			v := r.token(t, data)
			l := next()
			return &a.Cons{Car: v, Cdr: l}
		}
	}

	if data {
		return next()
	}
	return r.prepareList(next())
}

func (r *tokenReader) vector(data bool) a.Vector {
	var v = make(a.Vector, 0)

	for {
		t := r.lexer.Next()
		switch t.Type {
		case VectorEnd:
			return v
		case EndOfFile:
			panic(VectorNotClosed)
		default:
			e := r.token(t, data)
			v = append(v, e)
		}
	}
}

// EvalReader evaluates each element of the provided Reader
func EvalReader(c a.Context, reader Reader) a.Value {
	var r a.Value
	for v := reader.Next(); v != EndOfReader; v = reader.Next() {
		r = a.Eval(c, v)
	}
	return r
}
