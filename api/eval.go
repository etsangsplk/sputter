package api

// Evaluable identifies a Value as being directly evaluable
type Evaluable interface {
	Eval(Context) Value
}

// Eval is a ValueProcessor that expands and evaluates a Value
func Eval(c Context, v Value) Value {
	ex, _ := MacroExpand(c, v)
	if e, ok := ex.(Evaluable); ok {
		return e.Eval(c)
	}
	return ex
}
