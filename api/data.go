package api

// Quote identifies a Value as being in data mode (quoted)
type Quote struct {
	Value Value
}

// Eval makes Quote Evaluable
func (q *Quote) Eval(c Context) Value {
	return q.Value
}

func (q *Quote) String() string {
	return String(q.Value)
}
