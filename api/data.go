package api

// Data identifies a Value as being in data mode (literal)
type Data struct {
	Value Value
}

// Eval makes Data Evaluable
func (l *Data) Eval(c *Context) Value {
	return l.Value
}

func (l *Data) String() string {
	return String(l.Value)
}
