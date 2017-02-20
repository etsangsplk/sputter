package api

// UnknownSymbol is thrown if a symbol cannot be resolved
const UnknownSymbol = "symbol has not been defined"

// Symbol is an Identifier that can be resolved
type Symbol struct {
	Name Name
}

// Eval makes a Symbol Evaluable
func (s *Symbol) Eval(c Context) Value {
	if r, ok := c.Get(s.Name); ok {
		return r
	}
	panic(UnknownSymbol)
}

func (s *Symbol) String() string {
	return string(s.Name)
}
