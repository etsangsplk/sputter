package evaluator

import a "github.com/kode4food/sputter/api"

// ExpectedSequence is raised when unquote-splicing expands a non sequence
const ExpectedSequence = "can not Splice a non-sequence: %s"

// Splice contains Values that will be spliced into a macro expansion
type Splice []a.Value

// NewSplice converts a Value into a Splice if it is a Sequence
func NewSplice(v a.Value) Splice {
	if s, ok := v.(a.Sequence); ok {
		sp := Splice{}
		for i := s; i.IsSequence(); i = i.Rest() {
			sp = append(sp, i.First())
		}
		return sp
	}
	panic(a.Err(ExpectedSequence, v))
}

// Eval self-evaluates
func (s Splice) Eval(_ a.Context) a.Value {
	return s
}

// Str converts a Splice into a Str
func (s Splice) Str() a.Str {
	return a.MakeSequenceStr(a.NewVector(s...))
}
