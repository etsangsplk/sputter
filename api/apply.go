package api

// Applicable is the standard signature for any Value that can be applied
// to a sequence of arguments
type Applicable interface {
	Apply(Context, Vector) Value
}

// Apply applies a Sequence to an Applicable
func Apply(c Context, a Applicable, s Sequence) Value {
	return a.Apply(c, SequenceToVector(s))
}
