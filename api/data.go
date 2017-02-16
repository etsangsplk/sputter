package api

// Data identifies a Value as being in data mode (literal)
type Data struct {
	Value Value
}

// Evaluate makes Data Evaluable
func (l *Data) Eval(c *Context) Value {
	return l.Value
}

func (l *Data) String() string {
	return ValueToString(l.Value)
}
