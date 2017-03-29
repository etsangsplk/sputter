package api

// MetaMacro is the Metadata key that identifies a Function as being a Macro
var MetaMacro = NewKeyword("macro")

// Macro is a Function that can be used for Reader transformation
type Macro interface {
	Function
	DataMode() bool
}

type basicMacro struct {
	Function
	dataMode bool
}

var defaultMacroMetadata = Metadata{
	MetaMacro: true,
}

// NewMacro instantiates a new Macro
func NewMacro(e SequenceProcessor) Macro {
	f := NewFunction(e).WithMetadata(defaultMacroMetadata).(Function)
	return &basicMacro{
		Function: f,
		dataMode: true,
	}
}

func (m *basicMacro) DataMode() bool {
	return m.dataMode
}

// WithMetadata copies the Function with new Metadata
func (m *basicMacro) WithMetadata(md Metadata) Annotated {
	return &basicMacro{
		Function: m.Function.WithMetadata(md).(Function),
	}
}
