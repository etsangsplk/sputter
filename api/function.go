package api

// SequenceProcessor is the standard signature for a function that is
// capable of transforming or validating a Sequence
type SequenceProcessor func(Context, Sequence) Value

// Function is a Value that can be invoked
type Function struct {
	Name    Name
	Exec    SequenceProcessor
	Prepare SequenceProcessor
}

func (f *Function) String() string {
	return string(f.Name)
}
