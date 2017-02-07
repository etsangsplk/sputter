package interpreter

import (
	a "github.com/kode4food/sputter/api"
)

// ListNotClosed is the error returned when EOF is reached inside a List
const ListNotClosed = "end of file reached with open list"

// EndOfCoder represents the end of a Coder stream
var EndOfCoder = struct{}{}

// Coder is responsible for taking a stream of Tokens and converting them
// into Lists for evaluation
type Coder struct {
	builtIns *a.Context
	reader   TokenReader
}

// NewCoder instantiates a new Coder using the provided TokenReader
func NewCoder(builtIns *a.Context, reader TokenReader) *Coder {
	return &Coder{builtIns, reader}
}

// Next returns the next value from the Coder
func (c *Coder) Next() a.Value {
	return c.token(c.reader.Next())
}

func (c *Coder) token(t *Token) a.Value {
	switch t.Type {
	case LiteralMarker:
		return c.literal()
	case ListStart:
		return c.list(ListEnd)
	case ArgsStart:
		return c.list(ArgsEnd)
	case Identifier:
		return &a.Symbol{Name: t.Value.(string)}
	case EndOfFile:
		return EndOfCoder
	default:
		return t.Value
	}
}

func (c *Coder) literal() *a.Literal {
	return &a.Literal{Value: c.Next()}
}

func (c *Coder) list(endToken TokenType) *a.List {
	var handle func(token *Token) *a.List
	var first func() *a.List
	var next func() *a.List

	handle = func(token *Token) *a.List {
		switch token.Type {
		case endToken:
			return a.EmptyList
		case EndOfFile:
			panic(ListNotClosed)
		default:
			elem := c.token(token)
			list := next()
			return list.Cons(elem)
		}
	}

	first = func() *a.List {
		token := c.reader.Next()
		if token.Type == Identifier {
			name := token.Value.(string)
			if function, ok := c.builtIns.Get(name); ok {
				list := next()
				return list.Cons(function)
			}
		}
		return handle(token)
	}

	next = func() *a.List {
		token := c.reader.Next()
		return handle(token)
	}

	return first()
}
