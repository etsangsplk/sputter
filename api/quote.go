package api

// Quoted identifies a Value as being in data mode (quoted)
type Quoted interface {
	Value() Value
}

type quoted struct {
	value Value
}

// Quote wraps a Value into a data-mode container
func Quote(v Value) Quoted {
	return &quoted{
		value: v,
	}
}

// Value returns the Value that is Quoted
func (q quoted) Value() Value {
	return q.value
}

// Eval makes Quoted Evaluable
func (q quoted) Eval(c Context) Value {
	return q.value
}

func (q quoted) String() string {
	return String(q.value)
}
