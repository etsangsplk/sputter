package api

// SequenceProcessor is the standard signature for a function that is
// capable of transforming or validating a Sequence
type SequenceProcessor func(Context, Sequence) Value

// Function is a Value that can be invoked
type Function struct {
	Name    Name
	Exec    SequenceProcessor
	Prepare SequenceProcessor
	Data    bool
}

// ResolveFunction either returns a Function as-is or tries to perform
// a lookup if the Value is a Symbol
func ResolveFunction(c Context, v Value) (*Function, bool) {
	if f, ok := v.(*Function); ok {
		return f, true
	}

	if s, ok := v.(*Symbol); ok {
		if sv, ok := c.Get(s.Name); ok {
			if f, ok := sv.(*Function); ok {
				return f, true
			}
		}
	}

	return nil, false
}

func (f *Function) String() string {
	return string(f.Name)
}
