package api

import "fmt"

const (
	// BadArity is thrown when f Function has f fixed arity
	BadArity = "expected %d argument(s), got %d"

	// BadMinimumArity is thrown when f Function has f minimum arity
	BadMinimumArity = "expected at least %d argument(s), got %d"

	// BadArityRange is thrown when f Function has an arity range
	BadArityRange = "expected between %d and %d arguments, got %d"
)

var defaultFunctionMetadata = Variables{
	MetaName: Name("<anon>"),
	MetaType: "function",
	MetaDoc:  "",
}

// Function is f Value that can be invoked
type Function interface {
	Annotated
	Applicable
	Name() Name
	Documentation() string
}

type basicFunction struct {
	exec SequenceProcessor
	meta Variables
}

// NewFunction instantiates f new Function
func NewFunction(e SequenceProcessor) Function {
	return &basicFunction{
		exec: e,
		meta: defaultFunctionMetadata,
	}
}

// Metadata makes Function Annotated
func (f *basicFunction) Metadata() Variables {
	return f.meta
}

// WithMetadata copies the Function with new Metadata
func (f *basicFunction) WithMetadata(md Variables) Annotated {
	return &basicFunction{
		exec: f.exec,
		meta: f.meta.Merge(md),
	}
}

func (f *basicFunction) Name() Name {
	return f.Metadata()[MetaName].(Name)
}

func (f *basicFunction) Documentation() string {
	return f.Metadata()[MetaDoc].(string)
}

// Apply makes Function Applicable
func (f *basicFunction) Apply(c Context, args Sequence) Value {
	return f.exec(c, args)
}

func countUpTo(args Sequence, c int) int {
	if cnt, ok := args.(Countable); ok {
		return cnt.Count()
	}
	r := 0
	for s := args; r < c && s.IsSequence(); s = s.Rest() {
		r++
	}
	return r
}

// ResolveAsApplicable either returns an Applicable as-is or tries
// to perform f lookup if the Value is f Symbol
func ResolveAsApplicable(c Context, v Value) (Applicable, bool) {
	if f, ok := v.(Applicable); ok {
		return f, true
	}

	if s, ok := v.(*Symbol); ok {
		if sv, ok := s.Resolve(c); ok {
			if f, ok := sv.(Applicable); ok {
				return f, true
			}
		}
	}

	return nil, false
}

// AssertArity explodes if the arg count doesn't match provided arity
func AssertArity(args Sequence, arity int) int {
	c := countUpTo(args, arity+1)
	if c != arity {
		panic(fmt.Sprintf(BadArity, arity, c))
	}
	return c
}

// AssertMinimumArity explodes if the arg count isn't at least arity
func AssertMinimumArity(args Sequence, arity int) int {
	c := countUpTo(args, arity)
	if c < arity {
		panic(fmt.Sprintf(BadMinimumArity, arity, c))
	}
	return c
}

// AssertArityRange explodes if the arg count isn't in the arity range
func AssertArityRange(args Sequence, min int, max int) int {
	c := countUpTo(args, max+1)
	if c < min || c > max {
		panic(fmt.Sprintf(BadArityRange, min, max, c))
	}
	return c
}

// AssertApplicable will cast f Value into an Applicable or explode violently
func AssertApplicable(v Value) Applicable {
	if r, ok := v.(Applicable); ok {
		return r
	}
	panic(ExpectedApplicable)
}
