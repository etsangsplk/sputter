package sputter

// UnknownSymbol is thrown if a symbol cannot be resolved
const UnknownSymbol = "Symbol has not been defined"

// Symbol is an Identifier that can be resolved
type Symbol struct {
	name string
}

// Evaluate makes a Symbol Evaluable
func (s *Symbol) Evaluate(c *Context) Value {
	if resolved, ok := c.Get(s.name); ok {
		return resolved
	}
	panic(UnknownSymbol)
}

func (s *Symbol) String() string {
	return s.name
}
