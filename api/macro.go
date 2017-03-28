package api

// MetaMacro is the Metadata key that identifies a Function as being a Macro
const MetaMacro = "macro"

// Macro is a Function that can be used for Reader transformation
type Macro interface {
	Function
	DataMode() bool
}

type basicMacro struct {
	Function
	dataMode bool
}

var defaultMacroMetadata = Variables{
	MetaMacro: true,
}

// NewMacro instantiates f new Macro
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
func (m *basicMacro) WithMetadata(md Variables) Annotated {
	return &basicMacro{
		Function: m.Function.WithMetadata(md).(Function),
	}
}
