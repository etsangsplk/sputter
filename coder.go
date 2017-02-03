package sputter

// ListNotClosed is the error returned when EOF is reached inside a List
const ListNotClosed = "End of file reached with open list"

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
		return c.list()
	case Identifier:
		return &Resolvable{t.Value.(string)}
	default:
		return t.Value
	}
}

func (c *Coder) literal() *Literal {
	return &Literal{c.Next()}
}

func (c *Coder) list() *List {
	var handle func(token *Token) *List
	var first func() *List
	var next func() *List

	handle = func(token *Token) *List {
		switch token.Type {
		case ListEnd:
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
			list := next()
			return list.Cons(c.function(token))
		}
		return handle(token)
	}

	next = func() *List {
		token := c.reader.Next()
		return handle(token)
	}

	return first()
}

func (c *Coder) function(t *Token) *Function {
	var wrapper *Function
	name := t.Value.(string)

	wrapper = &Function{func(context *Context, list *List) Value {
		if v, f := context.Get(name); f {
			if entry, ok := v.(*Function); ok {
				// swap the exec function into the wrapper
				wrapper.exec = entry.exec
				return entry.exec(context, list)
			}
		}
		panic(NonFunction)
	}}

	return wrapper
}
