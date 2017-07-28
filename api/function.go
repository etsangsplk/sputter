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
	// ArityChecker is a function that validates the arity of arguments
	ArityChecker func(Sequence) (int, bool)

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

var (
	// DefaultFunctionName is the default name for anonymous functions
	DefaultFunctionName = Name("<lambda>")

	// DefaultFunctionType is the default type for functions
	DefaultFunctionType = Name("function")

	// MacroKey identifies a Function as being a Macro
	MacroKey = NewKeyword("macro")

	// SpecialKey identifies a Macro as being a special form
	SpecialKey = NewKeyword("special-form")

	functionMetadata = Properties{
		NameKey:    DefaultFunctionName,
		TypeKey:    DefaultFunctionType,
		MacroKey:   False,
		SpecialKey: False,
	}
)

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

func (f *function) WithMetadata(md Object) AnnotatedValue {
	return &function{
		exec: f.exec,
		meta: f.meta.Child(md.Flatten()),
	}
}

func (f *function) Name() Name {
	if v, ok := f.meta.Get(NameKey); ok {
		if n, ok := v.(Name); ok {
			return n
		}
	}
	return DefaultFunctionName
}

func (f *function) Documentation() Str {
	return GetDocumentation(f)
}

func (f *function) Type() Name {
	if v, ok := f.meta.Get(TypeKey); ok {
		if n, ok := v.(Name); ok {
			return n
		}
	}
	return DefaultFunctionType
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

// IsMacro tests an Applicable as being marked a Macro and is a special form
func IsMacro(a Applicable) (bool, bool) {
	if an, ok := a.(Annotated); ok {
		md := an.Metadata()
		return IsTrue(md, MacroKey), IsTrue(md, SpecialKey)
	}
	return false, false
}

// IsSpecialForm tests an Applicable as being marked a special form
func IsSpecialForm(a Applicable) bool {
	if an, ok := a.(Annotated); ok {
		return IsTrue(an.Metadata(), SpecialKey)
	}
	return false
}

// MakeArityChecker creates a fixed arity checker
func MakeArityChecker(arity int) ArityChecker {
	plusOne := arity + 1
	return func(args Sequence) (int, bool) {
		c := countUpTo(args, plusOne)
		return c, c == arity
	}
}

// AssertArity explodes if the arg count doesn't match provided arity
func AssertArity(args Sequence, arity int) int {
	c, ok := MakeArityChecker(arity)(args)
	if !ok {
		panic(ErrStr(BadArity, arity, c))
	}
	return c
}

// MakeMinimumArityChecker creates a minimum arity checker
func MakeMinimumArityChecker(arity int) ArityChecker {
	return func(args Sequence) (int, bool) {
		c := countUpTo(args, arity)
		return c, c >= arity
	}
}

// AssertMinimumArity explodes if the arg count isn't at least arity
func AssertMinimumArity(args Sequence, arity int) int {
	c, ok := MakeMinimumArityChecker(arity)(args)
	if !ok {
		panic(ErrStr(BadMinimumArity, arity, c))
	}
	return c
}

// MakeArityRangeChecker creates a ranged arity checker
func MakeArityRangeChecker(min int, max int) ArityChecker {
	maxPlusOne := max + 1
	return func(args Sequence) (int, bool) {
		c := countUpTo(args, maxPlusOne)
		return c, c >= min && c <= max
	}
}

// AssertArityRange explodes if the arg count isn't in the arity range
func AssertArityRange(args Sequence, min int, max int) int {
	c, ok := MakeArityRangeChecker(min, max)(args)
	if !ok {
		panic(ErrStr(BadArityRange, min, max, c))
	}
	return c
}
