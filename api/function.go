package api

import "reflect"

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

	baseFunction struct {
		meta Object
	}

	// HasReflectedFunction marks a Value as being based on reflection
	HasReflectedFunction interface {
		IsReflectedFunction() bool
	}

	// ReflectedFunction is the base structure for reflected Functions
	ReflectedFunction struct {
		baseFunction
		concrete reflect.Type
	}

	execFunction struct {
		baseFunction
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
)

// NewReflectedFunction uses reflection to instantiate a Base-derived Function
func NewReflectedFunction(f HasReflectedFunction) Function {
	t := reflect.ValueOf(f).Type().Elem()
	return newReflectedFunctionWithMeta(t, functionMetadata)
}

// NewExecFunction creates a Function instance from a SequenceProcessor
func NewExecFunction(e SequenceProcessor) Function {
	return &execFunction{
		baseFunction: baseFunction{
			meta: functionMetadata,
		},
		exec: e,
	}
}

func newReflectedFunctionWithMeta(t reflect.Type, md Object) Function {
	ptr := reflect.New(t)
	v := reflect.Indirect(ptr)
	f := reflect.Indirect(v).FieldByName("ReflectedFunction")
	f.Set(reflect.ValueOf(ReflectedFunction{
		baseFunction: baseFunction{
			meta: md,
		},
		concrete: t,
	}))
	return ptr.Interface().(Function)
}

func (f *baseFunction) IsFunction() bool {
	return true
}

func (f *baseFunction) Metadata() Object {
	return f.meta
}

func (f *baseFunction) Name() Name {
	if v, ok := f.Metadata().Get(NameKey); ok {
		if n, ok := v.(Name); ok {
			return n
		}
	}
	return DefaultFunctionName
}

func (f *baseFunction) Documentation() Str {
	return GetDocumentation(f.Metadata())
}

func (f *baseFunction) Type() Name {
	if v, ok := f.Metadata().Get(TypeKey); ok {
		if n, ok := v.(Name); ok {
			return n
		}
	}
	return DefaultFunctionType
}

func (f *baseFunction) Str() Str {
	return MakeDumpStr(f)
}

// IsReflectedFunction returns whether or not this Function is reflected
func (f *ReflectedFunction) IsReflectedFunction() bool {
	return true
}

// WithMetadata creates a copy of this Function with additional Metadata
func (f *ReflectedFunction) WithMetadata(md Object) AnnotatedValue {
	childMeta := f.Metadata().Child(md.Flatten())
	return newReflectedFunctionWithMeta(f.concrete, childMeta)
}

func (f *execFunction) WithMetadata(md Object) AnnotatedValue {
	return &execFunction{
		baseFunction: baseFunction{
			meta: f.meta.Child(md.Flatten()),
		},
		exec: f.exec,
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
