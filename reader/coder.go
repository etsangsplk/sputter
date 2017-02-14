package interpreter

import (
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
)

// EndOfCoder represents the end of a Coder stream
var EndOfCoder = struct{}{}

// Coder is responsible for taking a stream of Tokens and converting them
// into Lists for evaluation
type Coder struct {
	builtIns *a.Context
	reader   Reader
}

// NewCoder instantiates a new Coder using the provided Reader
func NewCoder(builtIns *a.Context, r Reader) *Coder {
	return &Coder{builtIns, r}
}

// Next returns the next value from the Coder
func (c *Coder) Next() a.Value {
	return c.token(c.reader.Next())
}

func (c *Coder) token(t *Token) a.Value {
	switch t.Type {
	case DataMarker:
		return c.data()
	case ListStart:
		return c.list()
	case VectorStart:
		return c.vector()
	case Identifier:
		return &a.Symbol{Name: a.Name(t.Value.(string))}
	case ListEnd:
		panic(UnmatchedListEnd)
	case VectorEnd:
		panic(UnmatchedVectorEnd)
	case EndOfFile:
		return EndOfCoder
	default:
		return t.Value
	}
}

func (c *Coder) data() *a.Data {
	return &a.Data{Value: c.Next()}
}

func (c *Coder) list() *a.List {
	var handle func(token *Token) *a.List
	var first func() *a.List
	var next func() *a.List

	handle = func(t *Token) *a.List {
		switch t.Type {
		case ListEnd:
			return a.EmptyList
		case EndOfFile:
			panic(ListNotClosed)
		default:
			v := c.token(t)
			l := next()
			return l.Cons(v)
		}
	}

	first = func() *a.List {
		t := c.reader.Next()
		if t.Type == Identifier {
			n := a.Name(t.Value.(string))
			if f, ok := c.builtIns.Get(n); ok {
				l := next()
				return l.Cons(f)
			}
		}
		return handle(t)
	}

	next = func() *a.List {
		t := c.reader.Next()
		return handle(t)
	}

	return first()
}

func (c *Coder) vector() a.Vector {
	var r = make(a.Vector, 0)

	for {
		t := c.reader.Next()
		switch t.Type {
		case VectorEnd:
			return r
		case EndOfFile:
			panic(VectorNotClosed)
		default:
			v := c.token(t)
			r = append(r, v)
		}
	}
}

// EvaluateCoder evaluates each element of the provided Coder
func EvaluateCoder(c *a.Context, coder *Coder) a.Value {
	var r a.Value
	for v := coder.Next(); v != EndOfCoder; v = coder.Next() {
		r = a.Evaluate(c, v)
	}
	return r
}
