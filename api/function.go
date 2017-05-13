package api

const (
	// BadArity is thrown when a Function has a fixed arity
	BadArity = "expected %d argument(s), got %d"

	// BadMinimumArity is thrown when a Function has a minimum arity
	BadMinimumArity = "expected at least %d argument(s), got %d"

	// BadArityRange is thrown when a Function has an arity range
	BadArityRange = "expected between %d and %d arguments, got %d"
)

var functionMetadata = Metadata{
	MetaName: Name("<anon>"),
	MetaType: Name("function"),
}

// Function is a Value that can be invoked
type Function struct {
	exec SequenceProcessor
	meta Metadata
}

// NewFunction instantiates a new Function
func NewFunction(e SequenceProcessor) *Function {
	return &Function{
		exec: e,
		meta: functionMetadata,
	}
}

// Metadata makes Function Annotated
func (f *Function) Metadata() Metadata {
	return f.meta
}

// WithMetadata copies the Function with new Metadata
func (f *Function) WithMetadata(md Metadata) Annotated {
	return &Function{
		exec: f.exec,
		meta: f.meta.Merge(md),
	}
}

// Name returns the name of this Function via its Metadata
func (f *Function) Name() Name {
	return f.Metadata()[MetaName].(Name)
}

// Documentation returns the docstring of this Function via its Metadata
func (f *Function) Documentation() Str {
	return f.Metadata()[MetaDoc].(Str)
}

// Type returns the Type of this Function via its Metadata
func (f *Function) Type() Name {
	if v, ok := f.meta.Get(MetaType); ok {
		if n, ok := v.(Name); ok {
			return n
		}
	}
	return "function"
}

// Apply makes Function Applicable
func (f *Function) Apply(c Context, args Sequence) Value {
	return f.exec(c, args)
}

// Str converts this Value into a Str
func (f *Function) Str() Str {
	return MakeDumpStr(f)
}

func countUpTo(args Sequence, c int) int {
	if cnt, ok := args.(Counted); ok {
		return cnt.Count()
	}
	r := 0
	for s := args; r < c && s.IsSequence(); s = s.Rest() {
		r++
	}
	return r
}

// AssertArity explodes if the arg count doesn't match provided arity
func AssertArity(args Sequence, arity int) int {
	c := countUpTo(args, arity+1)
	if c != arity {
		panic(Err(BadArity, arity, c))
	}
	return c
}

// AssertMinimumArity explodes if the arg count isn't at least arity
func AssertMinimumArity(args Sequence, arity int) int {
	c := countUpTo(args, arity)
	if c < arity {
		panic(Err(BadMinimumArity, arity, c))
	}
	return c
}

// AssertArityRange explodes if the arg count isn't in the arity range
func AssertArityRange(args Sequence, min int, max int) int {
	c := countUpTo(args, max+1)
	if c < min || c > max {
		panic(Err(BadArityRange, min, max, c))
	}
	return c
}
