package api

// Quoted identifies a Value as being in data mode (quoted)
type Quoted interface {
	Evaluable
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

// Eval makes Quoted Evaluable
func (q quoted) Eval(_ Context) Value {
	return q.value
}

func (q quoted) String() string {
	return String(q.value)
}
