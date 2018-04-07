package api

const (
	// DefaultFunctionName is the default name for anonymous functions
	DefaultFunctionName = Name("self")

	// DefaultFunctionType is the default type for functions
	DefaultFunctionType = Name("function")
)

type (
	// Function is a Value that can be invoked
	Function interface {
		Value
		Annotated
		Typed
		Documented
		Applicable
		FunctionType()
	}

	// BaseFunction provides common behavior for different Function types
	BaseFunction struct {
		meta Object
	}

	invokerFunction struct {
		BaseFunction
		invoke Invoker
	}

	blockFunction struct {
		BaseFunction
		body Block
	}
)

var (
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

// NewExecFunction creates a Function instance from a Invoker
func NewExecFunction(e Invoker) Function {
	return &invokerFunction{
		BaseFunction: DefaultBaseFunction,
		invoke:       e,
	}
}

// NewBlockFunction creates a new simple Function based on a Block
func NewBlockFunction(args Vector) Function {
	return &blockFunction{
		BaseFunction: DefaultBaseFunction,
		body:         MakeBlock(args),
	}
}

// Extend creates a new BaseFunction by extending its metadata
func (f *BaseFunction) Extend(md Object) BaseFunction {
	return BaseFunction{
		meta: f.meta.Child(md.Flatten()),
	}
}

// FunctionType identifies this Value as a Function
func (*BaseFunction) FunctionType() {}

// Metadata returns the Function's metadata Object
func (f *BaseFunction) Metadata() Object {
	return f.meta
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

func (f *invokerFunction) WithMetadata(md Object) AnnotatedValue {
	return &invokerFunction{
		BaseFunction: f.Extend(md),
		invoke:       f.invoke,
	}
}

func (f *invokerFunction) Apply(c Context, args Vector) Value {
	return f.invoke(c, args)
}

func (f *blockFunction) Apply(c Context, _ Vector) Value {
	return f.body.Eval(c)
}

func (f *blockFunction) WithMetadata(md Object) AnnotatedValue {
	return &blockFunction{
		BaseFunction: f.Extend(md),
		body:         f.body,
	}
}

// IsMacro tests an Applicable as being marked a Macro
func IsMacro(a Applicable) bool {
	if an, ok := a.(Annotated); ok {
		md := an.Metadata()
		return IsTrue(md, MacroKey)
	}
	return false
}

// IsSpecialForm tests an Applicable as being marked a special form
func IsSpecialForm(a Applicable) bool {
	if an, ok := a.(Annotated); ok {
		return IsTrue(an.Metadata(), SpecialKey)
	}
	return false
}
