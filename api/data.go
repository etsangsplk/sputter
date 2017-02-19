package api

// Atom identifies a Value as being atomic (evaluating to itself)
type Atom struct {
	Label string
}

// Eval makes Atom Evaluable
func (a *Atom) Eval(c *Context) Value {
	return a
}

func (a *Atom) String() string {
	return a.Label
}

// Quote identifies a Value as being in data mode (quoted)
type Quote struct {
	Value Value
}

// Eval makes Quote Evaluable
func (q *Quote) Eval(c *Context) Value {
	return q.Value
}

func (q *Quote) String() string {
	return String(q.Value)
}
