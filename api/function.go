package api

// ArgumentProcessor is the standard signature for a function that is
// capable of processing an Iterable (like Lists)
type ArgumentProcessor func(*Context, Iterable) Value

// Function is a Value that can be invoked
type Function struct {
	Name string
	Exec ArgumentProcessor
}

func (f *Function) String() string {
	return f.Name
}