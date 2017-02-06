package sputter

// ListNotClosed is the error returned when EOF is reached inside a List
const ListNotClosed = "End of file reached with open list"

// NonFunction is the error returned when a non-Function is invoked
const NonFunction = "first element of list is not a function"

// EndOfCoder represents the end of a Coder stream
var EndOfCoder = struct{}{}

// Coder is responsible for taking a stream of Tokens and converting them
// into Lists for evaluation
type Coder struct {
	reader TokenReader
}

// NewCoder instantiates a new Coder using the provided TokenReader
func NewCoder(reader TokenReader) *Coder {
	return &Coder{reader}
}

// Next returns the next value from the Coder
func (c *Coder) Next() Value {
	return c.token(c.reader.Next())
}

func (c *Coder) token(t *Token) Value {
	switch t.Type {
	case LiteralMarker:
		return c.literal()
	case ListStart:
		return c.list(ListEnd)
	case ArgsStart:
		return c.list(ArgsEnd)
	case Identifier:
		return &Symbol{t.Value.(string)}
	case EndOfFile:
		return EndOfCoder
	default:
		return t.Value
	}
}

func (c *Coder) literal() *Literal {
	return &Literal{c.Next()}
}

func (c *Coder) list(endToken TokenType) *List {
	var handle func(token *Token) *List
	var first func() *List
	var next func() *List

	handle = func(token *Token) *List {
		switch token.Type {
		case endToken:
			return EmptyList
		case EndOfFile:
			panic(ListNotClosed)
		default:
			elem := c.token(token)
			list := next()
			return list.Cons(elem)
		}
	}

	first = func() *List {
		token := c.reader.Next()
		if token.Type == Identifier {
			name := token.Value.(string)
			if function, ok := Builtins.Get(name); ok {
				list := next()
				return list.Cons(function)
			}
		}
		return handle(token)
	}

	next = func() *List {
		token := c.reader.Next()
		return handle(token)
	}

	return first()
}
