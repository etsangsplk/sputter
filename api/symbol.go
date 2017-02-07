package api

// UnknownSymbol is thrown if a symbol cannot be resolved
const UnknownSymbol = "Symbol has not been defined"

// Symbol is an Identifier that can be resolved
type Symbol struct {
	Name string
}

// Evaluate makes a Symbol Evaluable
func (s *Symbol) Evaluate(c *Context) Value {
	if resolved, ok := c.Get(s.Name); ok {
		return resolved
	}
	panic(UnknownSymbol)
}

func (s *Symbol) String() string {
	return s.Name
}
