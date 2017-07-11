package api

const (
	// BadArity is thrown when a Function has a fixed arity
	BadArity = "expected %d argument(s), got %d"

	// BadMinimumArity is thrown when a Function has a minimum arity
	BadMinimumArity = "expected at least %d argument(s), got %d"

	// BadArityRange is thrown when a Function has an arity range
	BadArityRange = "expected between %d and %d arguments, got %d"
)

type (
	// Function is a Value that can be invoked
	Function interface {
		Value
		Annotated
		Named
		Typed
		Documented
		Applicable
		IsFunction() bool
	}

	function struct {
		exec SequenceProcessor
		meta Object
	}
)

var functionMetadata = Properties{
	MetaName: Name("<anon>"),
	MetaType: Name("function"),
}

// NewFunction instantiates a new Function
func NewFunction(e SequenceProcessor) Function {
	return &function{
		exec: e,
		meta: functionMetadata,
	}
}

func (f *function) IsFunction() bool {
	return true
}

func (f *function) Metadata() Object {
	return f.meta
}

func (f *function) WithMetadata(md Object) Annotated {
	return &function{
		exec: f.exec,
		meta: f.meta.Child(md.Flatten()),
	}
}

func (f *function) Name() Name {
	n, _ := f.Metadata().Get(MetaName)
	return n.(Name)
}

func (f *function) Documentation() Str {
	d, _ := f.Metadata().Get(MetaDoc)
	return d.(Str)
}

func (f *function) Type() Name {
	if v, ok := f.meta.Get(MetaType); ok {
		if n, ok := v.(Name); ok {
			return n
		}
	}
	return "function"
}

func (f *function) Apply(c Context, args Sequence) Value {
	return f.exec(c, args)
}

func (f *function) Str() Str {
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
