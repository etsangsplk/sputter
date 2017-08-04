package api

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

	// BaseFunction provides common behavior for different Function types
	BaseFunction struct {
		meta Object
	}

	execFunction struct {
		BaseFunction
		exec SequenceProcessor
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

	// DefaultBaseFunction provides the default BaseFunction data
	DefaultBaseFunction = NewBaseFunction(functionMetadata)
)

// NewBaseFunction instantiates a new BaseFunction with metadata
func NewBaseFunction(md Object) BaseFunction {
	return BaseFunction{
		meta: md,
	}
}

// NewExecFunction creates a Function instance from a SequenceProcessor
func NewExecFunction(e SequenceProcessor) Function {
	return &execFunction{
		BaseFunction: DefaultBaseFunction,
		exec:         e,
	}
}

// Extend creates a new BaseFunction by extending its metadata
func (f *BaseFunction) Extend(md Object) BaseFunction {
	return BaseFunction{
		meta: f.meta.Child(md.Flatten()),
	}
}

// IsFunction identifies this Value as a Function
func (f *BaseFunction) IsFunction() bool {
	return true
}

// Metadata returns the Function's metadata Object
func (f *BaseFunction) Metadata() Object {
	return f.meta
}

// Name inspects the Function's metadata to determine its Name
func (f *BaseFunction) Name() Name {
	if v, ok := f.Metadata().Get(NameKey); ok {
		if n, ok := v.(Name); ok {
			return n
		}
	}
	return DefaultFunctionName
}

// Documentation returns the Function's documentation
func (f *BaseFunction) Documentation() Str {
	return GetDocumentation(f.Metadata())
}

// Type inspects the Function's metadata to determine its Type
func (f *BaseFunction) Type() Name {
	if v, ok := f.Metadata().Get(TypeKey); ok {
		if n, ok := v.(Name); ok {
			return n
		}
	}
	return DefaultFunctionType
}

// Str converts this Value to a Str
func (f *BaseFunction) Str() Str {
	return MakeDumpStr(f)
}

func (f *execFunction) WithMetadata(md Object) AnnotatedValue {
	return &execFunction{
		BaseFunction: f.Extend(md),
		exec:         f.exec,
	}
}

func (f *execFunction) Apply(c Context, args Sequence) Value {
	return f.exec(c, args)
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
